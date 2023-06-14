
#include <iostream>
#include <gpiod.hpp>
#include <unistd.h>
 
// https://stackoverflow.com/a/73186894

int main(int argc, char* argv[])
{ 
   ::gpiod::chip chip("gpiochip0");
   
   auto line = chip.get_line(17);  // GPIO17
   line.request({"example", gpiod::line_request::DIRECTION_OUTPUT, 0},1);  
   
   sleep(0.1);
   
   line.set_value(0);
   line.release();

   return EXIT_SUCCESS;
}
