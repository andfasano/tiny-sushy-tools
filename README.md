# tiny-sushy-tools

Very small (and not complete) Go porting of the original [sushy-tools project](https://github.com/openstack/sushy-tools) meant to support testing of the [Metal3 BareMetal Operator](https://github.com/metal3-io/baremetal-operator/)

# Pre-requisites

Install the libvirt-dev package:

```sudo yum install libvirt-devel libvirt-daemon-kvm libvirt-client```

(then start libvritd if required: ```sudo systemctl enable --now libvirtd```)