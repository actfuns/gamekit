package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/actfuns/gamekit/behavior_tree"
	"github.com/actfuns/gamekit/behavior_tree/actions"
	"github.com/actfuns/gamekit/behavior_tree/controls"
	"github.com/actfuns/gamekit/behavior_tree/decorators"
)

func main() {
	// 解析命令行参数
	testType := flag.String("test", "simple", "Test type: simple, complex, or full")
	flag.Parse()

	// 获取当前可执行文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	var xmlFile string
	switch *testType {
	case "simple":
		xmlFile = filepath.Join(dir, "test_tree.xml")
	case "complex":
		xmlFile = filepath.Join(dir, "complex_test.xml")
	case "full":
		xmlFile = filepath.Join(dir, "full_test.xml")
	default:
		log.Fatalf("Unknown test type: %s. Use 'simple', 'complex', or 'full'", *testType)
	}

	// 创建行为树工厂
	behaviorTreeFactory := behavior_tree.NewBehaviorTreeFactory()

	// 注册内置动作节点
	registerNodes(behaviorTreeFactory)

	// 创建XML解析器
	parser := behavior_tree.NewXMLParser(behaviorTreeFactory)

	tree, err := parser.LoadFromFile(xmlFile)
	if err != nil {
		log.Fatalf("Failed to load behavior tree from %s: %v", xmlFile, err)
	}

	fmt.Println("=== Behavior Tree XML Parser Demo ===")
	fmt.Printf("Successfully loaded behavior tree from: %s\n", xmlFile)

	// 打印树结构
	fmt.Println("\nBehavior Tree Structure:")
	tree.PrintTree()

	// 执行行为树
	fmt.Println("\nExecuting behavior tree...")
	status := tree.Tick()
	fmt.Printf("Behavior tree execution result: %s\n", status)

	// 演示多次执行
	fmt.Println("\nExecuting behavior tree multiple times:")
	for i := 0; i < 3; i++ {
		status = tree.Tick()
		fmt.Printf("Tick %d: %s\n", i+1, status)
	}

	// 展示如何使用黑板
	fmt.Println("\nBlackboard usage example:")
	blackboard := tree.Blackboard()
	blackboard.Set("test_value", "Hello World")
	blackboard.Set("counter", 42)

	if testValue, exists := blackboard.Get("test_value"); exists {
		fmt.Printf("Test value: %v\n", testValue)
	}
	if counter, exists := blackboard.Get("counter"); exists {
		fmt.Printf("Counter: %v\n", counter)
	}

	fmt.Println("\nDemo completed successfully!")
	os.Exit(0)
}

func registerNodes(factory *behavior_tree.BehaviorTreeFactory) {
	// 注册 AlwaysSuccess 节点
	err := factory.RegisterSimpleNode("AlwaysSuccess", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return actions.NewAlwaysSuccessNode(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register AlwaysSuccess: %v", err)
	}
	
	// 注册 AlwaysFailure 节点
	err = factory.RegisterSimpleNode("AlwaysFailure", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return actions.NewAlwaysFailureNode(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register AlwaysFailure: %v", err)
	}
	
	// 注册 Inverter 节点
	err = factory.RegisterSimpleNode("Inverter", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return decorators.NewInverterNode(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register Inverter: %v", err)
	}
	
	// 注册 Fallback 节点
	err = factory.RegisterSimpleNode("Fallback", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return controls.NewFallbackNode(name, config, false), nil
		})
	if err != nil {
		log.Fatalf("Failed to register Fallback: %v", err)
	}
	
	// 注册 ReactiveFallback 节点
	err = factory.RegisterSimpleNode("ReactiveFallback", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return controls.NewReactiveFallbackNode(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register ReactiveFallback: %v", err)
	}
	
	// 注册 Sequence 节点
	err = factory.RegisterSimpleNode("Sequence", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return controls.NewSequenceNode(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register Sequence: %v", err)
	}
	
	// 注册 ReactiveSequence 节点
	err = factory.RegisterSimpleNode("ReactiveSequence", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return controls.NewReactiveSequence(name, config), nil
		})
	if err != nil {
		log.Fatalf("Failed to register ReactiveSequence: %v", err)
	}
	
	// 注册 Sleep 节点（如果需要）
	err = factory.RegisterSimpleNode("Sleep", 
		func(name string, config behavior_tree.NodeConfig) (behavior_tree.TreeNode, error) {
			return &actions.SleepNode{}, nil
		})
	if err != nil {
		log.Fatalf("Failed to register Sleep: %v", err)
	}
}
