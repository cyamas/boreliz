package hold

type HoldType string

const (
	CRIMP  = "CRIMP"
	JUG    = "JUG"
	PINCH  = "PINCH"
	SLOPER = "SLOPER"
	JIB    = "JIB"
)

type Hold struct {
	id           int
	manufacturer string
	model        string
	color        string
	Type         HoldType
	wallID       int
	x, y         int
	angle        int
	width        int
	depth        float64
	incut        int
	texture      int
}

func New() *Hold {
	return &Hold{}
}

func (h *Hold) SetWallID(id int) {
	h.wallID = id
}

func (h *Hold) SetCoordinates(x, y int) {
	h.x, h.y = x, y
}

func (h *Hold) GetCoordinates() (x, y int) {
	return h.x, h.y
}

func (h *Hold) ID() int {
	return h.id
}

func (h *Hold) SetManufacturer(manufacturer string) {
	h.manufacturer = manufacturer
}

func (h *Hold) GetManufacturer() *string {
	return &h.manufacturer
}

func (h *Hold) SetModel(model string) {
	h.model = model
}

func (h *Hold) GetModel() string {
	return h.model
}

func (h *Hold) SetColor(color string) {
	h.color = color
}

func (h *Hold) GetColor() string {
	return h.color
}

func (h *Hold) SetAngle(angle int) {
	h.angle = angle
}

func (h *Hold) GetAngle() int {
	return h.angle
}

func (h *Hold) SetWidth(width int) {
	h.width = width
}

func (h *Hold) GetWidth() int {
	return h.width
}

func (h *Hold) SetDepth(depth float64) {
	h.depth = depth
}

func (h *Hold) GetDepth() float64 {
	return h.depth
}

func (h *Hold) SetIncut(incut int) {
	h.incut = incut
}

func (h *Hold) GetIncut() int {
	return h.incut
}

func (h *Hold) SetTexture(texture int) {
	h.texture = texture
}

func (h *Hold) GetTexture() int {
	return h.texture
}
