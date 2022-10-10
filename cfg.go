package xcfg

import (
	"embed"

	"github.com/gookit/goutil/fsutil"
	"gopkg.in/ini.v1"
)

//go:embed */*.cfg
var _efs embed.FS

const (
	embedCfgFile  = "raw/preset.cfg"
	customCfgFile = "~/.xcfg/xcfg.cfg"
)

type Ctrl struct {
	WithColor bool `ini:"with_color"`

	Debug    bool `ini:"debug"`
	DummyLog bool `ini:"dummy_log"`
	Recover  bool `ini:"recover"`
	Retry    uint `ini:"retry"`
	// should program pause in panic
	PauseInPanic bool `ini:"pause_in_panic"`

	// panic by mode
	PanicBy uint `ini:"panic_by,omitempty"`
}

type Log struct {
	Level  int  `ini:"level"`
	Caller bool `ini:"caller"`
	// with full caller name or not
	DefaultCaller bool `ini:"default_caller"`

	LogToConsole bool `ini:"log_to_console"`
	LogToFile    bool `ini:"log_to_file"`
	AsJson       bool `ini:"as_json"`
	MaxSize      int  `ini:"max_size"`
	MaxBackups   int  `ini:"max_backups"`
	MaxAge       int  `ini:"max_age"`

	SaveToDir string `ini:"save_to_dir"`
	FileName  string `ini:"file_name"`
}

type Updater struct {
	Provider string `ini:"provider"`

	URI  string `ini:"uri"`
	Name string `ini:"name"`
}

type PresetCfg struct {
	Ctrl Ctrl `ini:"ctrl"`
	Log  Log  `ini:"log"`

	Updater Updater `ini:"updater"`
}

func EfsRead(efs embed.FS, name string) []byte {
	raw, err := efs.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return raw
}

// LoadSources loads one or more data source
//
//	and the loading order
//	 1. embed: raw/preset.cfg
//	 2. customized: ~/.xcfg/xcfg.cfg
//	 3. raw byte passed in: WithRawData
//	 3. file passed in: WithConfigFile
//
// and in most cases, you can directly call MustLoadConfigFile
//
//	@return cfg
//	@return err
func LoadSources(opts ...OptFunc) (cfg *ini.File, err error) {
	opt := Opts{}
	bindOpts(&opt, opts...)

	cfg = ini.Empty()
	// 1. embed file preset.cfg
	dft := EfsRead(_efs, embedCfgFile)
	// 2. custom file ~/.xcfg/xcfg.cfg
	if fsutil.PathExist(customCfgFile) {
		custom := fsutil.ExpandPath(customCfgFile)
		err = cfg.Append(dft, custom)
		if err != nil {
			return
		}
	}

	// 3. rawDataList
	for _, v := range opt.rawList {
		err = cfg.Append(v)
		if err != nil {
			return
		}
	}

	// 4. custom files
	for _, f := range opt.files {
		f = fsutil.ExpandPath(f)
		err = cfg.Append(f)
		if err != nil {
			return
		}
	}

	cfg.BlockMode = false

	return cfg, err
}

func MustLoadSources(opts ...OptFunc) *ini.File {
	c, e := LoadSources(opts...)
	if e != nil {
		panic(e)
	}
	return c
}

func MustMapToCfg(ini *ini.File, cfg interface{}) {
	if err := ini.MapTo(cfg); err != nil {
		panic(err)
	}
}

var XCFG = &PresetCfg{}

func init() {
	var raw = MustLoadSources()
	MustMapToCfg(raw, XCFG)
}
