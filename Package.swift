// swift-tools-version: 6.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "binarytime",
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "binarytime",
            targets: ["binarytime"]
        ),
    ],
    targets: [
        // Targets are the basic building blocks of a package, defining a module or a test suite.
        // Targets can depend on other targets in this package and products from dependencies.
        .target(
            name: "binarytime"
        ),
        .testTarget(
            name: "binarytimeTests",
            dependencies: ["binarytime"]
        ),
    ],
    swiftLanguageModes: [.v6]
)
