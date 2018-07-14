package emulator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/Ruenzuo/nana/emulator"
)

var _ = Describe("Emulator", func() {
	var (
		emulator Emulator
	)

	BeforeEach(func() {
		emulator = *NewEmulator()
	})

	Describe("verifying operation codes work", func() {
		Context("when values overflow", func() {
			BeforeEach(func() {
				emulator.ProgramCounter.SetValue(0x0)
				emulator.StackPointer.SetValue(0xFFFF)
				emulator.ROM[0] = 0x0F
				emulator.ExecuteOpCode(0xF8)
			})

			It("should set the right flags", func() {
				Expect(emulator.FlagC()).To(Equal(true))
				Expect(emulator.FlagH()).To(Equal(true))
				Expect(emulator.FlagN()).To(Equal(false))
				Expect(emulator.FlagZ()).To(Equal(false))
			})
		})

		Context("when values don't overflow", func() {
			BeforeEach(func() {
				emulator.ProgramCounter.SetValue(0x0)
				emulator.StackPointer.SetValue(0xFFF0)
				emulator.ROM[0] = 0x01
				emulator.ExecuteOpCode(0xF8)
			})

			It("should set the right flags", func() {
				Expect(emulator.FlagC()).To(Equal(false))
				Expect(emulator.FlagH()).To(Equal(false))
				Expect(emulator.FlagN()).To(Equal(false))
				Expect(emulator.FlagZ()).To(Equal(false))
			})
		})
	})
})