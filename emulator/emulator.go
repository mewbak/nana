package emulator

import (
	"io/ioutil"
	"os"
)

const MaxCyclesPerSecond = 4194304
const MaxCyclesPerEmulationCycle = MaxCyclesPerSecond / 60 // Target is 60 FPS

type Emulator struct {
	AF Register16Bit
	BC Register16Bit
	DE Register16Bit
	HL Register16Bit

	CartridgeMemory [0x200000]uint8
	ROM             [0x10000]uint8
	RAM             [0x8000]uint8
	ProgramCounter  Register16Bit
	StackPointer    Register16Bit
	CurrentROMBank  uint16
	CurrentRAMBank  uint16
	Halted          bool
}

func NewEmulator() *Emulator {
	e := new(Emulator)
	e.ProgramCounter.SetValue(0x100)
	e.AF.SetValue(0x01B0)
	e.BC.SetValue(0x0013)
	e.DE.SetValue(0x00D8)
	e.HL.SetValue(0x014D)
	e.StackPointer.SetValue(0xFFFE)
	e.CurrentROMBank = 1 // Should never be 1, ROM bank 0 is fixed
	e.CurrentRAMBank = 0
	e.Halted = false
	e.ROM[0xFF05] = 0x00
	e.ROM[0xFF06] = 0x00
	e.ROM[0xFF07] = 0x00
	e.ROM[0xFF10] = 0x80
	e.ROM[0xFF11] = 0xBF
	e.ROM[0xFF12] = 0xF3
	e.ROM[0xFF14] = 0xBF
	e.ROM[0xFF16] = 0x3F
	e.ROM[0xFF17] = 0x00
	e.ROM[0xFF19] = 0xBF
	e.ROM[0xFF1A] = 0x7F
	e.ROM[0xFF1B] = 0xFF
	e.ROM[0xFF1C] = 0x9F
	e.ROM[0xFF1E] = 0xBF
	e.ROM[0xFF20] = 0xFF
	e.ROM[0xFF21] = 0x00
	e.ROM[0xFF22] = 0x00
	e.ROM[0xFF23] = 0xBF
	e.ROM[0xFF24] = 0x77
	e.ROM[0xFF25] = 0xF3
	e.ROM[0xFF26] = 0xF1
	e.ROM[0xFF40] = 0x91
	e.ROM[0xFF42] = 0x00
	e.ROM[0xFF43] = 0x00
	e.ROM[0xFF45] = 0x00
	e.ROM[0xFF47] = 0xFC
	e.ROM[0xFF48] = 0xFF
	e.ROM[0xFF49] = 0xFF
	e.ROM[0xFF4A] = 0x00
	e.ROM[0xFF4B] = 0x00
	e.ROM[0xFFFF] = 0x00
	return e
}

func (e *Emulator) LoadCartridge(filename string) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		os.Exit(1)
	}

	copy(e.CartridgeMemory[:], dat)
}

func (e *Emulator) EmulateSecond() {
	cyclesThisUpdate := 0
	for cyclesThisUpdate < MaxCyclesPerEmulationCycle {
		cycles := e.executeNextOpcode()
		cyclesThisUpdate += cycles
		e.executeInterrupts()
	}
}

func (e *Emulator) executeNextOpcode() int {
	opCode := e.ReadMemory8Bit(e.ProgramCounter.Value())
	e.ProgramCounter.Increment()
	var cycles int
	if e.Halted {
		cycles = 4
	} else {
		cycles = e.ExecuteOpCode(opCode)
	}
	return cycles
}

func (e *Emulator) executeInterrupts() {
	// TODO: move this after the verifications of interrupt enabled/requested
	e.Halted = false
}

func testBit(n uint8, pos uint) bool {
	mask := uint8(1) << pos
	return n&mask > 0
}

func setBit(n uint8, pos uint) uint8 {
	mask := uint8(1) << pos
	n |= mask
	return n
}

func clearBit(n uint8, pos uint) uint8 {
	mask := ^(uint8(1) << pos)
	n &= mask
	return n
}
