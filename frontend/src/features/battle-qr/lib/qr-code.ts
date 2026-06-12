const qrSize = 21
const dataCodewords = 19
const errorCodewords = 7
const formatErrorCorrectionLevel = 1
const maskPattern = 0
const maxByteLength = 17

type Module = boolean | null

export type QrMatrix = boolean[][]

export function createQrMatrix(value: string): QrMatrix {
  const data = createDataCodewords(value)
  const errorCorrection = reedSolomonRemainder(data, errorCodewords)
  const codewords = [...data, ...errorCorrection]
  const modules = createModules()
  const reserved = createReserved()

  const setFunction = (row: number, col: number, dark: boolean) => {
    if (row < 0 || row >= qrSize || col < 0 || col >= qrSize) return
    modules[row][col] = dark
    reserved[row][col] = true
  }

  drawFinderPattern(setFunction, 0, 0)
  drawFinderPattern(setFunction, 0, qrSize - 7)
  drawFinderPattern(setFunction, qrSize - 7, 0)
  drawTimingPatterns(setFunction)
  drawFormatBits(setFunction)
  setFunction(qrSize - 8, 8, true)
  drawDataModules(modules, reserved, codewords)

  return modules.map((row) => row.map(Boolean))
}

function createModules(): Module[][] {
  return Array.from({ length: qrSize }, () => Array<Module>(qrSize).fill(null))
}

function createReserved(): boolean[][] {
  return Array.from({ length: qrSize }, () =>
    Array<boolean>(qrSize).fill(false),
  )
}

function createDataCodewords(value: string) {
  const bytes = new TextEncoder().encode(value)
  if (bytes.length > maxByteLength) {
    throw new Error("QR code value is too long")
  }

  const bits: boolean[] = []
  appendBits(bits, 0b0100, 4)
  appendBits(bits, bytes.length, 8)
  bytes.forEach((byte) => appendBits(bits, byte, 8))

  const capacity = dataCodewords * 8
  appendBits(bits, 0, Math.min(4, capacity - bits.length))
  while (bits.length % 8 !== 0) {
    bits.push(false)
  }

  const codewords: number[] = []
  for (let index = 0; index < bits.length; index += 8) {
    codewords.push(bitsToByte(bits.slice(index, index + 8)))
  }

  for (
    let padIndex = 0;
    codewords.length < dataCodewords;
    padIndex = (padIndex + 1) % 2
  ) {
    codewords.push(padIndex === 0 ? 0xec : 0x11)
  }

  return codewords
}

function appendBits(bits: boolean[], value: number, length: number) {
  for (let index = length - 1; index >= 0; index -= 1) {
    bits.push(((value >>> index) & 1) === 1)
  }
}

function bitsToByte(bits: boolean[]) {
  return bits.reduce((value, bit) => (value << 1) | (bit ? 1 : 0), 0)
}

function drawFinderPattern(
  setFunction: (row: number, col: number, dark: boolean) => void,
  row: number,
  col: number,
) {
  for (let y = -1; y <= 7; y += 1) {
    for (let x = -1; x <= 7; x += 1) {
      const currentRow = row + y
      const currentCol = col + x
      const inFinder = x >= 0 && x <= 6 && y >= 0 && y <= 6
      const isOuter = x === 0 || x === 6 || y === 0 || y === 6
      const isCenter = x >= 2 && x <= 4 && y >= 2 && y <= 4

      setFunction(currentRow, currentCol, inFinder && (isOuter || isCenter))
    }
  }
}

function drawTimingPatterns(
  setFunction: (row: number, col: number, dark: boolean) => void,
) {
  for (let index = 8; index < qrSize - 8; index += 1) {
    const dark = index % 2 === 0
    setFunction(6, index, dark)
    setFunction(index, 6, dark)
  }
}

function drawFormatBits(
  setFunction: (row: number, col: number, dark: boolean) => void,
) {
  const bits = calculateFormatBits()

  for (let index = 0; index <= 5; index += 1) {
    setFunction(8, index, getBit(bits, index))
  }
  setFunction(8, 7, getBit(bits, 6))
  setFunction(8, 8, getBit(bits, 7))
  setFunction(7, 8, getBit(bits, 8))
  for (let index = 9; index < 15; index += 1) {
    setFunction(14 - index, 8, getBit(bits, index))
  }

  for (let index = 0; index < 8; index += 1) {
    setFunction(8, qrSize - 1 - index, getBit(bits, index))
  }
  for (let index = 8; index < 15; index += 1) {
    setFunction(qrSize - 15 + index, 8, getBit(bits, index))
  }
}

function calculateFormatBits() {
  const data = (formatErrorCorrectionLevel << 3) | maskPattern
  let remainder = data

  for (let index = 0; index < 10; index += 1) {
    remainder = (remainder << 1) ^ ((remainder >>> 9) * 0x537)
  }

  return ((data << 10) | remainder) ^ 0x5412
}

function getBit(value: number, index: number) {
  return ((value >>> index) & 1) !== 0
}

function drawDataModules(
  modules: Module[][],
  reserved: boolean[][],
  codewords: number[],
) {
  const bits = codewords.flatMap((codeword) =>
    Array.from({ length: 8 }, (_, index) => getBit(codeword, 7 - index)),
  )
  let bitIndex = 0
  let upward = true

  for (let right = qrSize - 1; right >= 1; right -= 2) {
    if (right === 6) right -= 1

    for (let vertical = 0; vertical < qrSize; vertical += 1) {
      const row = upward ? qrSize - 1 - vertical : vertical
      for (let col = right; col >= right - 1; col -= 1) {
        if (reserved[row][col]) continue

        let dark = bits[bitIndex] ?? false
        bitIndex += 1
        if ((row + col) % 2 === 0) {
          dark = !dark
        }
        modules[row][col] = dark
      }
    }

    upward = !upward
  }
}

const gfExp = new Array<number>(512)
const gfLog = new Array<number>(256)

let gfValue = 1
for (let index = 0; index < 255; index += 1) {
  gfExp[index] = gfValue
  gfLog[gfValue] = index
  gfValue <<= 1
  if (gfValue & 0x100) {
    gfValue ^= 0x11d
  }
}
for (let index = 255; index < gfExp.length; index += 1) {
  gfExp[index] = gfExp[index - 255]
}

function gfMultiply(left: number, right: number) {
  if (left === 0 || right === 0) return 0
  return gfExp[gfLog[left] + gfLog[right]]
}

function reedSolomonGenerator(degree: number) {
  let polynomial = [1]

  for (let degreeIndex = 0; degreeIndex < degree; degreeIndex += 1) {
    const next = Array<number>(polynomial.length + 1).fill(0)
    polynomial.forEach((coefficient, index) => {
      next[index] ^= coefficient
      next[index + 1] ^= gfMultiply(coefficient, gfExp[degreeIndex])
    })
    polynomial = next
  }

  return polynomial
}

function reedSolomonRemainder(data: number[], degree: number) {
  const generator = reedSolomonGenerator(degree)
  const result = Array<number>(degree).fill(0)

  data.forEach((byte) => {
    const factor = byte ^ result.shift()!
    result.push(0)
    for (let index = 0; index < degree; index += 1) {
      result[index] ^= gfMultiply(generator[index + 1], factor)
    }
  })

  return result
}
