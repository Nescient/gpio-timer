cmake_minimum_required(VERSION 3.14)

include(cmake/prelude.cmake)

project(
    gpio
    VERSION 0.1.0
    DESCRIPTION "A GPIO timer for DebyNet server"
    HOMEPAGE_URL "https://github.com/Nescient/gpio-timer"
    LANGUAGES CXX
)

include(cmake/project-is-top-level.cmake)
include(cmake/variables.cmake)

# ---- Declare library ----

add_library(
    gpio_lib OBJECT
    source/lib.cpp
)

target_include_directories(
    gpio_lib ${warning_guard}
    PUBLIC
    "$<BUILD_INTERFACE:${PROJECT_SOURCE_DIR}/source>"
)

target_compile_features(gpio_lib PUBLIC cxx_std_17)

find_package(fmt REQUIRED)
target_link_libraries(gpio_lib PRIVATE fmt::fmt)

# ---- Declare executable ----

add_executable(gpio_exe source/main.cpp)
add_executable(gpio::exe ALIAS gpio_exe)

set_property(TARGET gpio_exe PROPERTY OUTPUT_NAME gpio)

target_compile_features(gpio_exe PRIVATE cxx_std_17)

target_link_libraries(gpio_exe PRIVATE gpio_lib)

# ---- Install rules ----

if(NOT CMAKE_SKIP_INSTALL_RULES)
  include(cmake/install-rules.cmake)
endif()

# ---- Developer mode ----

if(NOT gpio_DEVELOPER_MODE)
  return()
elseif(NOT PROJECT_IS_TOP_LEVEL)
  message(
      AUTHOR_WARNING
      "Developer mode is intended for developers of gpio"
  )
endif()

include(cmake/dev-mode.cmake)
