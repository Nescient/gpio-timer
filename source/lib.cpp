#include "lib.hpp"

#include <fmt/core.h>
#include <iostream>

#include "restclient-cpp/restclient.h"

library::library()
    : name {fmt::format("{}", "gpio")}
{
    RestClient::Response r = RestClient::get("http://example.com");
    std::cout << r.body << std::endl;
}
