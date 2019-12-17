package oglconsts

// Buffer Bits
// nolint: golint,megacheck
const (
	DEPTH_BUFFER_BIT   uint32 = 0x00000100
	STENCIL_BUFFER_BIT        = 0x00000400
	COLOR_BUFFER_BIT          = 0x00004000
)

// Draw Types
// nolint: golint,megacheck
const (
	POINTS         uint32 = 0x0000
	LINE                  = 0x1B01
	LINES                 = 0x0001
	LINE_LOOP             = 0x0002
	LINE_STRIP            = 0x0003
	TRIANGLES             = 0x0004
	TRIANGLE_STRIP        = 0x0005
	TRIANGLE_FAN          = 0x0006
)

// Shader Types
// nolint: golint,megacheck
const (
	FRAGMENT_SHADER uint32 = 0x8B30
	VERTEX_SHADER          = 0x8B31
)

// Status Values
// nolint: golint,megacheck
const (
	COMPILE_STATUS  uint32 = 0x8B81
	LINK_STATUS            = 0x8B82
	INFO_LOG_LENGTH        = 0x8B84
	TRUE                   = 1
	FALSE                  = 0
)

// Buffer Types
// nolint: golint,megacheck
const (
	ARRAY_BUFFER         uint32 = 0x8892
	ELEMENT_ARRAY_BUFFER        = 0x8893
)

// Draw Types
// nolint: golint,megacheck
const (
	STREAM_DRAW  uint32 = 0x88E0
	STATIC_DRAW         = 0x88E4
	DYNAMIC_DRAW        = 0x88E8
)

// Features
// nolint: golint,megacheck
const (
	BLEND        uint32 = 0x0BE2
	DEPTH_TEST          = 0x0B71
	CULL_FACE           = 0x0B44
	SCISSOR_TEST        = 0x0C11
	LESS                = 0x0201
	LEQUAL              = 0x0203
	DEPTH_FUNC          = 0x0B74

	ACTIVE_TEXTURE               = 0x84E0
	CURRENT_PROGRAM              = 0x8B8D
	TEXTURE_BINDING_2D           = 0x8069
	SAMPLER_BINDING              = 0x8919
	ARRAY_BUFFER_BINDING         = 0x8894
	ELEMENT_ARRAY_BUFFER_BINDING = 0x8895
	VERTEX_ARRAY_BINDING         = 0x85B5
	POLYGON_MODE                 = 0x0B40
	VIEWPORT                     = 0x0BA2
	SCISSOR_BOX                  = 0x0C10
	BLEND_SRC_RGB                = 0x80C9
	BLEND_DST_RGB                = 0x80C8
	BLEND_SRC_ALPHA              = 0x80CB
	BLEND_DST_ALPHA              = 0x80CA
	BLEND_EQUATION_RGB           = 0x8009
	BLEND_EQUATION_ALPHA         = 0x883D
	FRONT_AND_BACK               = 0x0408
	FILL                         = 0x1B02
)

// Alpha constants
// nolint: golint,megacheck
const (
	SRC_ALPHA           uint32 = 0x0302
	ONE_MINUS_SRC_ALPHA        = 0x0303
	ONE_MINUS_SRC_COLOR        = 0x0301

	FUNC_ADD = 0x8006
)

// Data Types
// nolint: golint,megacheck
const (
	BYTE           uint32 = 0x1400
	UNSIGNED_BYTE         = 0x1401
	SHORT                 = 0x1402
	UNSIGNED_SHORT        = 0x1403
	INT                   = 0x1404
	UNSIGNED_INT          = 0x1405
	FLOAT                 = 0x1406
)

// Texture Constants
// nolint: golint,megacheck
const (
	TEXTURE_2D uint32 = 0x0DE1

	TEXTURE0 = 0x84C0

	TEXTURE_CUBE_MAP            = 0x8513
	TEXTURE_CUBE_MAP_POSITIVE_X = 0x8515
	TEXTURE_CUBE_MAP_NEGATIVE_X = 0x8516
	TEXTURE_CUBE_MAP_POSITIVE_Y = 0x8517
	TEXTURE_CUBE_MAP_NEGATIVE_Y = 0x8518
	TEXTURE_CUBE_MAP_POSITIVE_Z = 0x8519
	TEXTURE_CUBE_MAP_NEGATIVE_Z = 0x851A

	NEAREST            = 0x2600
	TEXTURE_MAG_FILTER = 0x2800
	TEXTURE_MIN_FILTER = 0x2801
	TEXTURE_WRAP_S     = 0x2802
	TEXTURE_WRAP_T     = 0x2803
	TEXTURE_WRAP_R     = 0x8072
	CLAMP_TO_EDGE      = 0x812F

	UNPACK_ROW_LENGTH = 0x0CF2

	LINEAR = 0x2601
)

// Errors
// nolint: golint,megacheck
const (
	NO_ERROR                      uint32 = 0
	INVALID_ENUM                         = 0x0500
	INVALID_VALUE                        = 0x0501
	INVALID_OPERATION                    = 0x0502
	STACK_OVERFLOW                       = 0x0503
	STACK_UNDERFLOW                      = 0x0504
	OUT_OF_MEMORY                        = 0x0505
	INVALID_FRAMEBUFFER_OPERATION        = 0x0506
)

// Color Types
// nolint: golint,megacheck
const (
	ALPHA uint32 = 0x1906
	RGBA         = 0x1908
	RED          = 0x1903
	R8           = 0x8229
)
