package xcfg

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/ungerik/go-dry"
)

type CfgSuite struct {
	suite.Suite
}

func TestCfg(t *testing.T) {
	suite.Run(t, new(CfgSuite))
}

func (s *CfgSuite) SetupSuite() {
}

func (s *CfgSuite) TearDownSuite() {
}

func (s *CfgSuite) Test_01_LoadConfig() {
	raw := `
[ctrl]
pause_in_panic = true
`
	file0 := "/tmp/xcfg_cfg_00.cfg"
	dry.FileSetString(file0, raw)

	raw1 := `
[ctrl]
pause_in_panic = false
`
	file1 := "/tmp/xcfg_cfg_01.cfg"
	dry.FileSetString(file1, raw1)

	file2 := "/tmp/xcfg_cfg_01.cfg"

	cfg := &PresetCfg{}
	ini := MustLoadSources(WithConfigFile(file2), WithConfigFile(file0))
	MustMapToCfg(ini, cfg)

	s.IsType(Ctrl{}, cfg.Ctrl)
	s.True(cfg.Ctrl.PauseInPanic)

	ini = MustLoadSources(WithConfigFile(file0), WithConfigFile(file1))
	MustMapToCfg(ini, cfg)
	s.False(cfg.Ctrl.PauseInPanic)

	ini = MustLoadSources(WithConfigFile(file1), WithConfigFile(file0))
	MustMapToCfg(ini, cfg)
	s.True(cfg.Ctrl.PauseInPanic)
}

func (s *CfgSuite) Test_02_EfsRead() {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args

		wantByte  []byte
		wantPanic bool
	}{
		{
			name:      "non-exist should panics",
			args:      args{name: "non-exist"},
			wantByte:  []byte{},
			wantPanic: true,
		},
		{
			name: "preset.cfg",
			args: args{
				name: embedCfgFile,
			},
			wantByte:  []byte{},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		if tt.wantPanic {
			s.Panics(func() {
				EfsRead(_efs, tt.args.name)
			}, tt.name)
		} else {
			b := EfsRead(_efs, tt.args.name)
			s.NotNil(b)
		}
	}
}

func (s *CfgSuite) Test_03_CFG() {
	ini := MustLoadSources()
	cfg := &PresetCfg{}

	MustMapToCfg(ini, cfg)
	s.IsType(&PresetCfg{}, cfg)
	s.True(cfg.Ctrl.Debug, "debug should be true")

	s.Equal(XCFG, cfg)

	raw := `
[ctrl]
pause_in_panic = true
`
	file0 := "/tmp/xcfg_cfg_00.cfg"
	dry.FileSetString(file0, raw)

	ini = MustLoadSources(WithConfigFile(file0))
	MustMapToCfg(ini, cfg)

	s.NotEqual(XCFG, cfg)
}
