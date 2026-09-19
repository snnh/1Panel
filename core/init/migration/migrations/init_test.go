package migrations

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/dto"
)

func TestRemoveXpackAndAiMenus(t *testing.T) {
	menus := []dto.ShowMenu{
		{ID: "1", Label: "Home-Menu", Path: "/"},
		{ID: "20", Label: "Xpack-Menu", Path: "/xpack", Children: []dto.ShowMenu{
			{ID: "201", Label: "Node", Path: "/xpack/node"},
			{ID: "202", Label: "Cluster", Path: "/xpack/cluster"},
		}},
		{ID: "30", Label: "AI-Menu", Path: "/ai", Children: []dto.ShowMenu{
			{ID: "301", Label: "Agents", Path: "/ai/agents/agent"},
		}},
		{ID: "40", Label: "XApp", Path: "/xpack/app"},
		{ID: "50", Label: "System-Menu", Path: "/hosts/files", Children: []dto.ShowMenu{
			{ID: "501", Label: "File", Path: "/hosts/files"},
			{ID: "502", Label: "Sync", Path: "/xpack/sync/file"},
			{ID: "503", Label: "AIModel", Path: "/ai/model/account"},
		}},
	}

	kept := removeXpackAndAiMenus(menus)

	labels := make([]string, 0, len(kept))
	for _, menu := range kept {
		labels = append(labels, menu.Label)
	}
	expected := []string{"Home-Menu", "System-Menu"}
	if len(labels) != len(expected) {
		t.Fatalf("expected menus %v, got %v", expected, labels)
	}
	for i := range expected {
		if labels[i] != expected[i] {
			t.Fatalf("expected menus %v, got %v", expected, labels)
		}
	}

	var children []dto.ShowMenu
	for _, menu := range kept {
		if menu.Label == "System-Menu" {
			children = menu.Children
		}
	}
	if len(children) != 1 || children[0].Path != "/hosts/files" {
		t.Fatalf("expected only the local file menu, got %+v", children)
	}
}
