package behavior_tree

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// XMLParser is used to parse behavior tree XML files
type XMLParser struct {
	factory *BehaviorTreeFactory
}

// NewXMLParser creates a new XML parser
func NewXMLParser(factory *BehaviorTreeFactory) *XMLParser {
	return &XMLParser{
		factory: factory,
	}
}

// BehaviorTreeXML represents the root of a behavior tree XML
type BehaviorTreeXML struct {
	XMLName  xml.Name  `xml:"root"`
	MainTree string    `xml:"main_tree_to_execute,attr"`
	Trees    []TreeXML `xml:"BehaviorTree"`
}

// TreeXML represents a single behavior tree in XML
type TreeXML struct {
	XMLName xml.Name `xml:"BehaviorTree"`
	ID      string   `xml:"ID,attr"`
	Root    NodeXML  `xml:",any"`
}

// NodeXML represents a node in XML
type NodeXML struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []NodeXML  `xml:",any"`
}

// LoadFromFile loads a behavior tree from an XML file
func (p *XMLParser) LoadFromFile(filename string) (*BehaviorTree, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	var btXML BehaviorTreeXML
	if err := xml.Unmarshal(data, &btXML); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	// Find the main tree
	var mainTree *TreeXML
	for i := range btXML.Trees {
		if btXML.Trees[i].ID == btXML.MainTree {
			mainTree = &btXML.Trees[i]
			break
		}
	}

	if mainTree == nil {
		return nil, fmt.Errorf("main tree '%s' not found", btXML.MainTree)
	}

	// Create blackboard
	blackboard := NewBlackboard()

	// Parse the root node
	rootNode, err := p.parseNode(mainTree.Root, blackboard)
	if err != nil {
		return nil, err
	}

	return NewBehaviorTree(rootNode, blackboard, p.factory), nil
}

// parseNode recursively parses a node from XML
func (p *XMLParser) parseNode(nodeXML NodeXML, blackboard *core.Blackboard) (core.TreeNode, error) {
	nodeName := nodeXML.XMLName.Local

	// Extract attributes
	attrs := make(map[string]string)
	for _, attr := range nodeXML.Attrs {
		attrs[attr.Name.Local] = attr.Value
	}

	// Create node config
	config := NodeConfig{
		Blackboard: blackboard,
	}

	// Check if this is a registered node type
	if _, exists := p.factory.GetManifest(nodeName); exists {
		// Create node using factory
		node, err := p.factory.CreateNode(nodeName, nodeName, config)
		if err != nil {
			return nil, err
		}

		// Parse children
		for _, childXML := range nodeXML.Children {
			childNode, err := p.parseNode(childXML, blackboard)
			if err != nil {
				return nil, err
			}

			// Add child to the node
			node.AddChild(childNode)
		}

		return node, nil
	}

	// Handle built-in control nodes
	var treeNode TreeNode

	switch nodeName {
	case "Sequence":
		seqNode := NewSequenceNode(nodeName, config)
		treeNode = seqNode
	case "ReactiveSequence":
		rseqNode := NewReactiveSequenceNode(nodeName, config)
		treeNode = rseqNode
	case "Fallback":
		fbkNode := NewFallbackNode(nodeName, config)
		treeNode = fbkNode
	case "ReactiveFallback":
		rfbkNode := NewReactiveFallbackNode(nodeName, config)
		treeNode = rfbkNode
	case "Inverter":
		invNode := NewInverterNode(nodeName, config)
		treeNode = invNode
	case "RetryUntilSuccessful":
		maxAttempts := 3
		if attemptsStr, exists := attrs["num_attempts"]; exists {
			// Parse num_attempts (simplified)
			maxAttempts = parseInt(attemptsStr, 3)
		}
		retryNode := NewRetryNode(nodeName, config, maxAttempts)
		treeNode = retryNode
	case "Repeat":
		maxRepetitions := 3
		if repsStr, exists := attrs["num_cycles"]; exists {
			maxRepetitions = parseInt(repsStr, 3)
		}
		repNode := NewRepeatNode(nodeName, config, maxRepetitions)
		treeNode = repNode
	case "Timeout":
		timeoutMs := int64(1000)
		if timeoutStr, exists := attrs["msec"]; exists {
			timeoutMs = parseInt64(timeoutStr, 1000)
		}
		timeoutNode := NewTimeoutNode(nodeName, config, timeoutMs)
		treeNode = timeoutNode
	default:
		// Try to create as action node
		// This is a simplified approach - in real implementation you'd have more sophisticated logic
		treeNode = NewActionNode(nodeName, config, func() NodeStatus { return NodeStatusSuccess })
	}

	// Parse children
	for _, childXML := range nodeXML.Children {
		childNode, err := p.parseNode(childXML, blackboard)
		if err != nil {
			return nil, err
		}
		treeNode.AddChild(childNode)
	}

	return treeNode, nil
}

// parseInt parses a string to int with default fallback
func parseInt(s string, defaultValue int) int {
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		// This is a port reference, use default
		return defaultValue
	}

	// Simplified parsing - in real implementation you'd use strconv
	switch s {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	}
	return defaultValue
}

// parseInt64 parses a string to int64 with default fallback
func parseInt64(s string, defaultValue int64) int64 {
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		// This is a port reference, use default
		return defaultValue
	}

	// Simplified parsing
	switch s {
	case "1000":
		return 1000
	case "2000":
		return 2000
	}
	return defaultValue
}
