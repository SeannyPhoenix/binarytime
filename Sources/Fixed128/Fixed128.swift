import Foundation

enum Fixed128Error: Error {
    case divisionByZero
}

public struct Fixed128: Equatable {
    public let neg: Bool
    public let hi: UInt64
    public let lo: UInt64

    public init(neg: Bool, hi: UInt64, lo: UInt64) {
        self.neg = neg
        self.hi = hi
        self.lo = lo
    }
}

public func NewFixed128(x: Int64, y: Int64) throws -> Fixed128 {
    return try toFixed128(x: x, y: y)
}

func toFixed128(x: Int64, y: Int64) throws -> Fixed128 {
    if y == 0 {
        throw Fixed128Error.divisionByZero
    }

    let (negX, absX) = normalize(value: x)
    let (negY, absY) = normalize(value: y)
    let neg = negX != negY

    let (hi, lo) = getComponents(x: absX, y: absY)

    return Fixed128(neg: neg, hi: hi, lo: lo)
}

func normalize(value: Int64) -> (neg: Bool, abs: UInt64) {
    let mask: UInt64 = UInt64(value >> 63)
    let neg: Bool = mask != 0
    let abs: UInt64 = (UInt64(value) ^ mask) - mask
    return (neg, abs)
}

func getComponents(x: UInt64, y: UInt64) -> (UInt64, UInt64) {
    var hi: UInt64 = 0
    var lo: UInt64 = 0

    hi = x / y
    var part = x % y

    let shift = y.leadingZeroBitCount
    var div = y << shift
    part <<= shift

    var i = 0
    while i < 64 && div > 1 && part > 0 {
        div >>= 1
        let bit = 1 & ~UInt64(bitPattern: Int64(bitPattern: part - y) >> 63)
        part -= bit * div
        lo <<= 1
        lo |= bit
        i += 1
    }

    lo <<= (64 - i)
    return (hi, lo)
}
