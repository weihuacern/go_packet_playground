package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"plugin"

	"github.com/google/gopacket"

	phttp "plugin/http/plugin"
)

// PluginHandler : Plugin handler to load internal or external traffic handler
type PluginHandler struct {
	// Public
	ResolveStream      func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
	BPF                string
	InternalPluginList map[string]PluginInterface
	ExternalPluginList map[string]ExternalPlugin

	// Private
	dir string
}

// PluginInterface : A interface defined plugin
type PluginInterface interface {
	Version() string
	SetFlag([]string)
	BPFFilter() string
	ResolveStream(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
}

// ExternalPlugin : External plugin object
type ExternalPlugin struct {
	Name          string
	Version       string
	SetFlag       func([]string)
	BPFFilter     func() string
	ResolvePacket func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
}

// NewPluginHandler : Instantiate a new PluginHandler
func NewPluginHandler() *PluginHandler {
	var p PluginHandler
	p.dir, _ = filepath.Abs("./plug/")
	p.LoadInternalPluginList()
	p.LoadExternalPluginList()
	return &p
}

// LoadInternalPluginList : Load plugins from internal
func (p *PluginHandler) LoadInternalPluginList() {
	list := make(map[string]PluginInterface)
	// HTTP
	list["http"] = phttp.NewInstance()
	p.InternalPluginList = list
}

// LoadExternalPluginList : Load plugins from external
func (p *PluginHandler) LoadExternalPluginList() {
	dir, err := ioutil.ReadDir(p.dir)
	if err != nil {
		return
	}

	p.ExternalPluginList = make(map[string]ExternalPlugin)
	for _, fi := range dir {
		if fi.IsDir() || path.Ext(fi.Name()) != ".so" {
			continue
		}
		plug, err := plugin.Open(p.dir + "/" + fi.Name())
		if err != nil {
			panic(err)
		}
		VersionFunc, err := plug.Lookup("Version")
		if err != nil {
			panic(err)
		}
		Version := VersionFunc.(func() string)()
		SetFlagFunc, err := plug.Lookup("SetFlag")
		if err != nil {
			panic(err)
		}
		BPFFilterFunc, err := plug.Lookup("BPFFilter")
		if err != nil {
			panic(err)
		}
		ResolvePacketFunc, err := plug.Lookup("ResolvePacket")
		if err != nil {
			panic(err)
		}
		p.ExternalPluginList[fi.Name()] = ExternalPlugin{
			Name:          fi.Name(),
			Version:       Version,
			SetFlag:       SetFlagFunc.(func([]string)),
			BPFFilter:     BPFFilterFunc.(func() string),
			ResolvePacket: ResolvePacketFunc.(func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)),
		}
	}
}

// PrintPluginList : Print Plugins that have been loaded
func (p *PluginHandler) PrintPluginList() {
	// Print Internal Plugins
	for inPluginName := range p.InternalPluginList {
		fmt.Printf("Internal plugin: %s\n", inPluginName)
	}

	fmt.Printf("-- --- --\n")

	// Print External Plugin
	for exPluginName := range p.ExternalPluginList {
		fmt.Printf("External plugin: %s\n", exPluginName)
	}
}

// SetOption : Set options, like BPF, etc., to the plugin
func (p *PluginHandler) SetOption(pluginName string, pluginParams []string) {
	// Load Internal Plugin
	if internalPlugin, ok := p.InternalPluginList[pluginName]; ok {
		p.ResolveStream = internalPlugin.ResolveStream
		internalPlugin.SetFlag(pluginParams)
		p.BPF = internalPlugin.BPFFilter()
		return
	}

	// Load External Plugin
	externalPlugin, err := plugin.Open("./plug/" + pluginName)
	if err != nil {
		panic(err)
	}
	resolvePacket, err := externalPlugin.Lookup("ResolvePacket")
	if err != nil {
		panic(err)
	}
	setFlag, err := externalPlugin.Lookup("SetFlag")
	if err != nil {
		panic(err)
	}
	BPFFilter, err := externalPlugin.Lookup("BPFFilter")
	if err != nil {
		panic(err)
	}
	p.ResolveStream = resolvePacket.(func(net gopacket.Flow, transport gopacket.Flow, r io.Reader))
	setFlag.(func([]string))(pluginParams)
	p.BPF = BPFFilter.(func() string)()
}
