cmake_minimum_required(VERSION 3.5)

add_compile_options(-std=c++11)
add_compile_options(-Wall)
add_compile_options(-Wextra)
add_compile_options(-Werror)

add_definitions(-DSHADERS_DIR="shaders")

add_executable(scene src/main.cpp)

find_package(PkgConfig REQUIRED)

pkg_search_module(GLFW REQUIRED glfw3)
include_directories(${GLFW_INCLUDE_DIRS})
target_link_libraries(scene ${GLFW_LIBRARIES})

pkg_search_module(EPOXY REQUIRED epoxy)
include_directories(${EPOXY_INCLUDE_DIRS})
target_link_libraries(scene ${EPOXY_LIBRARIES})

include_directories(${PROJECT_SOURCE_DIR}/include)

find_package(Boost COMPONENTS system filesystem REQUIRED)
target_link_libraries(scene ${Boost_FILESYSTEM_LIBRARY} ${Boost_SYSTEM_LIBRARY})