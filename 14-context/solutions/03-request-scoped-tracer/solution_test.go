package tracer

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestTracer_AttachAndRetrieve(t *testing.T) {
	ctx := context.Background()

	// Initially absent
	if _, ok := GetTraceInfo(ctx); ok {
		t.Fatal("expected no trace info in background context")
	}

	info := TraceInfo{
		TraceID:   "trace-001",
		RequestID: "req-abc",
		UserID:    "user-42",
		Tags:      map[string]string{"service": "orders", "region": "us-east-1"},
	}

	ctxWithTrace := WithTraceInfo(ctx, info)

	extracted, ok := GetTraceInfo(ctxWithTrace)
	if !ok {
		t.Fatal("expected trace info to be present")
	}

	if extracted.TraceID != "trace-001" || extracted.RequestID != "req-abc" || extracted.UserID != "user-42" {
		t.Fatalf("unexpected trace info fields: %+v", extracted)
	}

	if extracted.Tags["service"] != "orders" || extracted.Tags["region"] != "us-east-1" {
		t.Fatalf("unexpected tags: %+v", extracted.Tags)
	}
}

func TestTracer_ImmutabilityAndBranching(t *testing.T) {
	ctx := context.Background()
	parentCtx := WithTraceInfo(ctx, TraceInfo{
		TraceID: "root-trace",
		Tags:    map[string]string{"env": "prod"},
	})

	childCtx1, err := AddTraceTag(parentCtx, "branch", "child-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	childCtx2, err := AddTraceTag(parentCtx, "branch", "child-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify parent has not been mutated
	parentInfo, _ := GetTraceInfo(parentCtx)
	if _, hasBranch := parentInfo.Tags["branch"]; hasBranch {
		t.Errorf("parent context tags should not contain 'branch'")
	}

	// Verify children have distinct branch tags
	child1Info, _ := GetTraceInfo(childCtx1)
	if child1Info.Tags["branch"] != "child-1" {
		t.Errorf("child1 branch tag mismatch: %s", child1Info.Tags["branch"])
	}

	child2Info, _ := GetTraceInfo(childCtx2)
	if child2Info.Tags["branch"] != "child-2" {
		t.Errorf("child2 branch tag mismatch: %s", child2Info.Tags["branch"])
	}
}

func TestTracer_ConcurrentAccessSafety(t *testing.T) {
	ctx := WithTraceInfo(context.Background(), TraceInfo{
		TraceID: "concurrent-trace",
		Tags:    map[string]string{"seed": "true"},
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			info, ok := GetTraceInfo(ctx)
			if !ok || info.TraceID != "concurrent-trace" {
				t.Errorf("worker %d failed to read trace info", workerID)
			}
			// Attempt to mutate the returned map - should not affect other workers
			info.Tags["mutated"] = "danger"
		}(i)
	}
	wg.Wait()

	finalInfo, _ := GetTraceInfo(ctx)
	if _, mutated := finalInfo.Tags["mutated"]; mutated {
		t.Fatalf("stored context tags were mutated concurrently!")
	}
}

func TestTracer_AddTagNoTraceInfoError(t *testing.T) {
	ctx := context.Background()
	_, err := AddTraceTag(ctx, "key", "value")
	if !errors.Is(err, ErrNoTraceInfo) {
		t.Fatalf("expected ErrNoTraceInfo, got %v", err)
	}
}
