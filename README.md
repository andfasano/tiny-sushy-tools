# tiny-sushy-tools

Very small (and not complete) Go porting of the original [sushy-tools project](https://github.com/openstack/sushy-tools) meant to support testing of the [Metal3 BareMetal Operator](https://github.com/metal3-io/baremetal-operator/)

## Building the package

- Pre-requisites for building the package: Install the libvirt-dev package

    ```bash
    sudo dnf install libvirt-devel libvirt-daemon-kvm libvirt-client
    ```

- Make sure `libvirt` service is running in the backend server

    ```bash
    sudo systemctl enable --now libvirtd
    ```

- Compiling the code

    ```bash
    go build -o tiny-sushy cmd/sushy-mock/main.go
    ```

- Running the code. (Pre-requisite for user to have ssh-key authentication to the libvirt node)

    ```bash
    ./tiny-sushy -ip 192.168.1.10 -port 9000
    2021/09/08 00:07:43 Starting RedFish mock server on port  9000
    ```

- `curl` queries for testing configuration

    ```bash
    curl http://localhost:9000/redfish/v1/Systems/<your-system-id>
    ```
