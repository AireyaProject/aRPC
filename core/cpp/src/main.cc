#include <iostream>
#include "aireya/server.h"
#include "../../../example.aya.h"

using namespace aireya;
using namespace aireya::example;

class UserServiceImpl : public UserService {
public:
    int GetUser(const User& req, User* resp) override {
        std::cout << "Received GetUser request for ID: " << req.id << std::endl;
        
        // Mock response
        resp->id = req.id;
        resp->name = "Aireya Pioneer";
        return 0; // Success
    }

    int StreamUsers(const User& req, User* resp) override {
        (void)req;
        (void)resp;
        // Not implemented in this basic example
        return -1;
    }
};

int main() {
    try {
        Server server(8080);
        UserServiceImpl service;

        server.SetHandler([&service](std::shared_ptr<Connection> conn, const Frame& frame) {
            std::cout << "Received Frame of type " << static_cast<int>(frame.type) 
                      << " with payload size " << frame.payload.size() << "\n";

            if (frame.type == FrameType::REQUEST) {
                User req, resp;
                if (req.ParseFromArray(frame.payload.data(), frame.payload.size())) {
                    
                    // Dispatch to service (In reality, we'd use metadata to route to the correct RPC method)
                    service.GetUser(req, &resp);

                    // Send back response
                    Frame resp_frame;
                    resp_frame.type = FrameType::RESPONSE;
                    resp_frame.stream_id = frame.stream_id;
                    resp.SerializeToArray(resp_frame.payload);
                    
                    conn->SendFrame(resp_frame);
                } else {
                    std::cerr << "Failed to parse User request payload.\n";
                }
            }
        });

        server.Run();
    } catch (const std::exception& e) {
        std::cerr << "Server failed: " << e.what() << "\n";
        return 1;
    }

    return 0;
}
