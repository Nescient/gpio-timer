#include "lib.hpp"

#include <fmt/core.h>
#include <iostream>


#include "restclient-cpp/restclient.h"

std::string BASE_URL{"localhost"};
std::string ACTION_URL{"/derbynet/action.php"};

library::library()
    : name {fmt::format("{}", "gpio")}
{
    RestClient::Response r = RestClient::get("http://example.com");
    std::cout << r.body << std::endl;

    // initialize RestClient
RestClient::init();

// get a connection object
connection = std::make_shared<RestClient::Connection>("localhost");

// configure basic auth
connection->SetBasicAuth("Timer", "");

// set connection timeout to 5s
connection->SetTimeout(5);

//action=timer-message&message=HEARTBEAT
RestClient::Response r2 = connection->post("/action.php", "application/x-www-form-urlencoded", R"({"action": "timer-message", "message": "HELLO"})");

std::cout << r2.body << std::endl;
}
