cmake_minimum_required(VERSION 3.14)
project(AireyaCore CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Compiler flags for high performance
add_compile_options(-O3 -Wall -Wextra -Werror)

# Include directories
include_directories(include)

# Source files
set(SOURCES
    src/wire_format.cc
    src/frame.cc
    src/server.cc
    src/quic_server.cc
    src/resolver.cc
    src/crypto.cc
)

# Build the core static library
add_library(aireya_core STATIC ${SOURCES})
target_include_directories(aireya_core PUBLIC include)

# Build the example server
add_executable(aireya_server src/main.cc)
target_include_directories(aireya_server PRIVATE ../..)
target_link_libraries(aireya_server PRIVATE aireya_core)
