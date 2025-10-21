package helper

import (
	"K8SArdoqBridge/app/controllers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// MockArdoqServer provides a mock HTTP server for Ardoq API calls
type MockArdoqServer struct {
	*httptest.Server
	components map[string]controllers.Component
	references map[string]controllers.Reference
	workspaces map[string]controllers.Workspace
	models     map[string]controllers.Model
	mu         sync.RWMutex
	RequestLog []MockRequest
}

type MockRequest struct {
	Method string
	Path   string
	Body   string
}

// NewMockArdoqServer creates a new mock Ardoq API server
func NewMockArdoqServer() *MockArdoqServer {
	mock := &MockArdoqServer{
		components: make(map[string]controllers.Component),
		references: make(map[string]controllers.Reference),
		workspaces: make(map[string]controllers.Workspace),
		models:     make(map[string]controllers.Model),
		RequestLog: make([]MockRequest, 0),
	}

	// Initialize default workspace and model
	mock.workspaces["test-workspace-id"] = controllers.Workspace{
		ID:             "test-workspace-id",
		Name:           "Test Workspace",
		ComponentModel: "test-model-id",
	}

	mock.models["test-model-id"] = controllers.Model{
		ID:          "test-model-id",
		Name:        "Test Model",
		Description: "Test Model Description",
		Root: controllers.ModelComponentTypes{
			"Cluster": {
				ID:   "cluster-type-id",
				Name: "Cluster",
			},
			"Namespace": {
				ID:   "namespace-type-id",
				Name: "Namespace",
			},
			"Deployment": {
				ID:   "deployment-type-id",
				Name: "Deployment",
			},
			"StatefulSet": {
				ID:   "statefulset-type-id",
				Name: "StatefulSet",
			},
			"Node": {
				ID:   "node-type-id",
				Name: "Node",
			},
			"SharedResourceComponent": {
				ID:   "shared-resource-type-id",
				Name: "SharedResourceComponent",
			},
			"SharedNodeComponent": {
				ID:   "shared-node-type-id",
				Name: "SharedNodeComponent",
			},
		},
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(mock.handler))
	return mock
}

func (m *MockArdoqServer) handler(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Read and save the body so we can log it and still pass it to handlers
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		if err := r.Body.Close(); err != nil {
			log.Errorf("Mock API server error: %s", err)
		}
		// Create a new reader with the saved bytes
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	m.RequestLog = append(m.RequestLog, MockRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   string(bodyBytes),
	})

	log.Tracef("Mock Ardoq API: %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	// Route the request
	switch {
	case strings.HasPrefix(r.URL.Path, "/workspace/"):
		m.handleWorkspace(w, r)
	case strings.HasPrefix(r.URL.Path, "/model/"):
		m.handleModel(w, r)
	case strings.HasPrefix(r.URL.Path, "/component/search"):
		m.handleComponentSearch(w, r)
	case strings.HasPrefix(r.URL.Path, "/component/"):
		m.handleComponent(w, r)
	case r.URL.Path == "/component":
		m.handleComponentCreate(w, r)
	case strings.HasPrefix(r.URL.Path, "/reference/"):
		m.handleReference(w, r)
	case r.URL.Path == "/reference":
		if r.Method == http.MethodGet {
			m.handleReferenceList(w, r)
		} else {
			m.handleReferenceCreate(w, r)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Not found"})
	}
}

func (m *MockArdoqServer) handleWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/workspace/")
	if workspace, ok := m.workspaces[id]; ok {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(workspace)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Workspace not found"})
	}
}

func (m *MockArdoqServer) handleModel(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/model/")
	if model, ok := m.models[id]; ok {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Model not found"})
	}
}

func (m *MockArdoqServer) handleComponentSearch(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspace")
	name := r.URL.Query().Get("name")

	results := make([]controllers.Component, 0)
	for id, comp := range m.components {
		if comp.RootWorkspace == workspace && (name == "" || comp.Name == name) {
			// Ensure fields map is initialized and has safe defaults
			if comp.Fields == nil {
				comp.Fields = make(map[string]interface{})
			}
			m.ensureSafeFields(comp.Fields, comp.Type)
			// CRITICAL: Store the component back with populated defaults
			m.components[id] = comp
			results = append(results, comp)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(results)
}

func (m *MockArdoqServer) handleComponent(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/component/")

	switch r.Method {
	case http.MethodGet:
		if comp, ok := m.components[id]; ok {
			// Ensure fields map is initialized and has safe defaults
			if comp.Fields == nil {
				comp.Fields = make(map[string]interface{})
			}
			m.ensureSafeFields(comp.Fields, comp.Type)
			// CRITICAL: Store the component back with populated defaults
			m.components[id] = comp
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(comp)
		} else {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Component not found"})
		}
	case http.MethodPut, http.MethodPatch:
		// Decode as a generic map since BodyProvider flattens the request
		var flatRequest map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Invalid request"})
			return
		}
		// Get existing component or create new one
		comp, exists := m.components[id]
		if !exists {
			comp = controllers.Component{ID: id}
		}
		// Update fields from request
		if name := getString(flatRequest, "name"); name != "" {
			comp.Name = name
		}
		if desc := getString(flatRequest, "description"); desc != "" {
			comp.Description = desc
		}
		if workspace := getString(flatRequest, "rootWorkspace"); workspace != "" {
			comp.RootWorkspace = workspace
		}
		if typeID := getString(flatRequest, "typeId"); typeID != "" {
			comp.TypeID = typeID
			// Update Type name when TypeID changes
			comp.Type = m.getTypeNameFromID(typeID)
		}
		if parent := getString(flatRequest, "parent"); parent != "" {
			comp.Parent = parent
		}
		comp.Fields = flatRequest // Store all fields
		// CRITICAL: Ensure safe defaults after update
		m.ensureSafeFields(comp.Fields, comp.Type)
		m.components[id] = comp
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(comp)
	case http.MethodDelete:
		if _, ok := m.components[id]; ok {
			delete(m.components, id)
			// Also delete associated references
			for refID, ref := range m.references {
				if ref.Source == id || ref.Target == id {
					delete(m.references, refID)
				}
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
		} else {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Component not found"})
		}
	}
}

func (m *MockArdoqServer) handleComponentCreate(w http.ResponseWriter, r *http.Request) {
	// Decode as a generic map since BodyProvider flattens the request
	var flatRequest map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Invalid request"})
		return
	}

	typeID := getString(flatRequest, "typeId")
	typeName := m.getTypeNameFromID(typeID)

	// Ensure safe defaults for fields based on component type
	m.ensureSafeFields(flatRequest, typeName)

	// Ensure Fields map is initialized
	if flatRequest == nil {
		flatRequest = make(map[string]interface{})
		m.ensureSafeFields(flatRequest, typeName)
	}

	// Build Component from flat request
	comp := controllers.Component{
		ID:            RandomString(24),
		Name:          getString(flatRequest, "name"),
		Description:   getString(flatRequest, "description"),
		RootWorkspace: getString(flatRequest, "rootWorkspace"),
		TypeID:        typeID,
		Type:          typeName, // Map TypeID to Type name
		Parent:        getString(flatRequest, "parent"),
		Fields:        flatRequest, // Store all fields
	}

	m.components[comp.ID] = comp

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comp)
}

// ensureSafeFields ensures that fields have safe default values to prevent panics
func (m *MockArdoqServer) ensureSafeFields(fields map[string]interface{}, typeName string) {
	// Helper to set default if key doesn't exist or is nil
	setDefault := func(key string, defaultVal interface{}) {
		if _, ok := fields[key]; !ok || fields[key] == nil {
			fields[key] = defaultVal
		}
	}

	// Helper to convert numeric fields to strings (matching Ardoq API behavior)
	ensureString := func(key string) {
		if val, ok := fields[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				fields[key] = fmt.Sprintf("%.0f", v)
			case int:
				fields[key] = fmt.Sprintf("%d", v)
			case int32:
				fields[key] = fmt.Sprintf("%d", v)
			case int64:
				fields[key] = fmt.Sprintf("%d", v)
			}
		}
	}

	switch typeName {
	case "Node":
		setDefault("node_capacity_cpu", float64(0))
		setDefault("node_capacity_memory", "0")
		setDefault("node_capacity_storage", "0")
		setDefault("node_capacity_pods", float64(0))
		setDefault("node_allocatable_cpu", float64(0))
		setDefault("node_allocatable_memory", "0")
		setDefault("node_allocatable_storage", "0")
		setDefault("node_allocatable_pods", float64(0))
		setDefault("node_container_runtime", "")
		setDefault("node_kernel_version", "")
		setDefault("node_kubelet_version", "")
		setDefault("node_kube_proxy_version", "")
		setDefault("node_os_image", "")
		setDefault("node_provider", "")
		setDefault("node_creation_timestamp", "")
	case "Deployment", "StatefulSet":
		setDefault("resource_image", "")
		setDefault("resource_replicas", "0")
		setDefault("resource_creation_timestamp", "")
		setDefault("resource_requests_cpu", "")
		setDefault("resource_requests_memory", "")
		setDefault("resource_limits_cpu", "")
		setDefault("resource_limits_memory", "")
		// Convert numeric replicas to string
		ensureString("resource_replicas")
	case "SharedResourceComponent":
		// Shared resource components might be accessed as resources
		setDefault("resource_image", "")
		setDefault("resource_replicas", "0")
		setDefault("resource_creation_timestamp", "")
		setDefault("resource_requests_cpu", "")
		setDefault("resource_requests_memory", "")
		setDefault("resource_limits_cpu", "")
		setDefault("resource_limits_memory", "")
		// Convert numeric replicas to string
		ensureString("resource_replicas")
	case "SharedNodeComponent":
		// Shared node components might be accessed as nodes
		setDefault("node_capacity_cpu", float64(0))
		setDefault("node_capacity_memory", "0")
		setDefault("node_capacity_storage", "0")
		setDefault("node_capacity_pods", float64(0))
		setDefault("node_allocatable_cpu", float64(0))
		setDefault("node_allocatable_memory", "0")
		setDefault("node_allocatable_storage", "0")
		setDefault("node_allocatable_pods", float64(0))
		setDefault("node_container_runtime", "")
		setDefault("node_kernel_version", "")
		setDefault("node_kubelet_version", "")
		setDefault("node_kube_proxy_version", "")
		setDefault("node_os_image", "")
		setDefault("node_provider", "")
		setDefault("node_creation_timestamp", "")
	}
}

// getTypeNameFromID maps a type ID to its name
func (m *MockArdoqServer) getTypeNameFromID(typeID string) string {
	for _, model := range m.models {
		for typeName, typeInfo := range model.Root {
			if typeInfo.ID == typeID {
				return typeName
			}
		}
	}
	return ""
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (m *MockArdoqServer) handleReference(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/reference/")

	switch r.Method {
	case http.MethodGet:
		if ref, ok := m.references[id]; ok {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ref)
		} else {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Reference not found"})
		}
	case http.MethodDelete:
		if _, ok := m.references[id]; ok {
			delete(m.references, id)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
		} else {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Reference not found"})
		}
	}
}

func (m *MockArdoqServer) handleReferenceList(w http.ResponseWriter, r *http.Request) {
	// Return all references as an array
	refs := make([]controllers.Reference, 0, len(m.references))
	for _, ref := range m.references {
		refs = append(refs, ref)
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(refs)
}

func (m *MockArdoqServer) handleReferenceCreate(w http.ResponseWriter, r *http.Request) {
	// Decode as a generic map since BodyProvider flattens the request
	var flatRequest map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(controllers.HttpError{Message: "Invalid request"})
		return
	}

	// Build Reference from flat request
	ref := controllers.Reference{
		ID:              RandomString(24),
		DisplayText:     getString(flatRequest, "displayText"),
		Description:     getString(flatRequest, "description"),
		RootWorkspace:   getString(flatRequest, "rootWorkspace"),
		TargetWorkspace: getString(flatRequest, "targetWorkspace"),
		Source:          getString(flatRequest, "source"),
		Target:          getString(flatRequest, "target"),
		Type:            getInt(flatRequest, "type"),
		Fields:          flatRequest, // Store all fields
	}

	m.references[ref.ID] = ref

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ref)
}

// Helper function to safely get int from map
func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return 0
}

// GetComponent retrieves a component by ID from the mock server
func (m *MockArdoqServer) GetComponent(id string) (controllers.Component, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	comp, ok := m.components[id]
	return comp, ok
}

// GetReference retrieves a reference by ID from the mock server
func (m *MockArdoqServer) GetReference(id string) (controllers.Reference, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ref, ok := m.references[id]
	return ref, ok
}

// ComponentCount returns the number of components in the mock server
func (m *MockArdoqServer) ComponentCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.components)
}

// ReferenceCount returns the number of references in the mock server
func (m *MockArdoqServer) ReferenceCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.references)
}

// Reset clears all stored data except the default workspace and model
func (m *MockArdoqServer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.components = make(map[string]controllers.Component)
	m.references = make(map[string]controllers.Reference)
	m.RequestLog = make([]MockRequest, 0)
}

// ConfigureMockEnvironment sets up environment variables to use this mock server
func (m *MockArdoqServer) ConfigureMockEnvironment() error {
	if err := os.Setenv("ARDOQ_BASEURI", m.URL); err != nil {
		return err
	}
	if err := os.Setenv("ARDOQ_APIKEY", "test-api-key"); err != nil {
		return err
	}
	if err := os.Setenv("ARDOQ_ORG", "test-org"); err != nil {
		return err
	}
	if err := os.Setenv("ARDOQ_WORKSPACE_ID", "test-workspace-id"); err != nil {
		return err
	}
	return nil
}
