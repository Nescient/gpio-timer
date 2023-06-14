#include <iostream>
#include <string>

#include "lib.hpp"

auto main(int argc, char* argv[]) -> int
{
  auto const lib = library {};
  auto const message = "Hello from " + lib.name + "!";
  std::cout << message << '\n';
  return EXIT_SUCCESS;
}
