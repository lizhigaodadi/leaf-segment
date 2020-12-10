package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/EslRain/leaf-segment/model"
	"sync/atomic"
	"time"
)

const (
	MAXSTEP  = 10e6
	Default  = 2000
	MAXRETRY = 3 // 最大重试次数
)

func (s *Service) CreateLeaf(ctx context.Context, leaf *model.Leaf) error {
	if leaf.Step > MAXSTEP {
		return errors.New("step limit exceeded")
	}

	if len(leaf.BizTag) == 0 {
		return errors.New("param invalid")
	}

	// 等于0 则代表使用默认step
	if leaf.Step == 0 {
		leaf.Step = Default
	}
	if leaf.MaxID == 0 {
		leaf.MaxID = 1
	}

	return s.dao.CreateLeaf(ctx, leaf)
}

func (s *Service) GetLeafCache(bizTag string) *model.LeafCache {
	if cache, ok := s.leafCache.Load(bizTag); ok {
		return cache.(*model.LeafCache)
	}
	return nil
}

func (s *Service) UpdateLeafCache(key string, bean *model.LeafCache) {
	if e, ok := s.leafCache.Load(key); ok {
		alloc := e.(*model.LeafCache)
		alloc.Buffer = bean.Buffer
		alloc.UpdateTime = bean.UpdateTime
	}
}

func (s *Service) PreloadBuffer(ctx context.Context, bizTag string, current *model.LeafCache) error {
	for i := 0; i < MAXRETRY; i++ {
		leaf, err := s.dao.NextSegment(ctx, bizTag)
		if err != nil {
			fmt.Printf("preloadBuffer failed; bizTag:%s;err:%v", bizTag, err)
			continue
		}

		segment := &model.LeafSegment{
			Cursor: leaf.MaxID - uint64(leaf.Step+1),
			Max:    leaf.MaxID - 1,
			Min:    leaf.MaxID - uint64(leaf.Step),
			InitOk: true,
		}
		current.Buffer = append(current.Buffer, segment)
		s.UpdateLeafCache(bizTag, current)
		current.Mutex.Lock()
		defer current.Mutex.Unlock()
		for _, waitChan := range current.Waiting {
			close(waitChan)
		}
		current.Waiting = current.Waiting[:0]
		break
	}

	current.IsPreload = false
	return nil
}

func (s *Service) GetID(ctx context.Context, bizTag string) (uint64, error) {
	leafCache := s.GetLeafCache(bizTag)

	//进行初始化缓存
	if leafCache == nil {
		leaf, err := s.dao.NextSegment(ctx, bizTag)
		if err != nil {
			fmt.Printf("initCache failed; err:%v\n", err)
			return 0, err
		}
		newLeafCache := &model.LeafCache{
			Key:        bizTag,
			Step:       leaf.Step,
			CurrentPos: 0,
			Buffer:     make([]*model.LeafSegment, 0),
			UpdateTime: time.Now(),
			Waiting:    make([]chan byte, 0),
			IsPreload:  false,
		}
		newLeafSegment := &model.LeafSegment{
			Cursor: leaf.MaxID - uint64(leaf.Step+1),
			Max:    leaf.MaxID - 1,
			Min:    leaf.MaxID - uint64(leaf.Step),
			InitOk: true,
		}
		newLeafCache.Buffer = append(newLeafCache.Buffer, newLeafSegment)
		s.leafCache.Store(bizTag, newLeafCache)
		leafCache = newLeafCache
	}

	leafCache.Mutex.Lock()
	defer leafCache.Mutex.Unlock()
	//首次获取id看下能否获取成功
	var id uint64
	currentBuffer := leafCache.Buffer[leafCache.CurrentPos]
	if currentBuffer.InitOk && currentBuffer.Cursor < currentBuffer.Max {
		id = atomic.AddUint64(&currentBuffer.Cursor, 1)
		leafCache.UpdateTime = time.Now()
	}

	//预加载
	if currentBuffer.Max-id < uint64(0.9*float32(leafCache.Step)) && len(leafCache.Buffer) <= 1 && !leafCache.IsPreload {
		leafCache.IsPreload = true
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		go s.PreloadBuffer(ctx, leafCache.Key, leafCache)
	}

	//如果一个buffer分发到最大值
	if id == currentBuffer.Max {
		if len(leafCache.Buffer) > 1 && leafCache.Buffer[leafCache.CurrentPos+1].InitOk {
			leafCache.Buffer = append(leafCache.Buffer[:0], leafCache.Buffer[1:]...)
		}
	}

	//如果分发到id了则直接返回
	if id > 0 {
		return id, nil
	}

	waitChan := make(chan byte, 1)
	leafCache.Waiting = append(leafCache.Waiting, waitChan)
	leafCache.Mutex.Unlock()

	timer := time.NewTimer(500 * time.Millisecond)
	select {
	case <-waitChan:
	case <-timer.C:
	}

	leafCache.Mutex.Lock()
	if len(leafCache.Buffer) <= 1 {
		return 0, errors.New("get id failed")
	}

	leafCache.Buffer = append(leafCache.Buffer[:0], leafCache.Buffer[1:]...)
	currentBuffer = leafCache.Buffer[leafCache.CurrentPos]
	if currentBuffer.InitOk && currentBuffer.Cursor < currentBuffer.Max {
		id = atomic.AddUint64(&leafCache.Buffer[leafCache.CurrentPos].Cursor, 1)
		leafCache.UpdateTime = time.Now()
	}
	return id, nil
}
