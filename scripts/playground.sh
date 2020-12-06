#!/bin/bash

# this network example
# A(.1)  --172.27.1.0/24--  (.2)B(.2)  --172.27.2.0/24-- C(.1)


set -eu

if [[ $(id -u) -ne 0 ]] ; then
    echo "Please run with sudo"
    exit 1
fi

run () {
    echo "$@"
    "$@" || exit 1
}

create_network () {
    __create_network
}

__create_network () {
    # setup namespaces
    run ip netns add HOSTA
    run ip netns add HOSTB
    run ip netns add HOSTC

    # setup veth peer
    run ip link add veth-A-B type veth peer name veth-B-A
    run ip link add veth-C-B type veth peer name veth-B-C

    run ip link set veth-A-B netns HOSTA
    run ip link set veth-B-A netns HOSTB
    run ip link set veth-B-C netns HOSTB
    run ip link set veth-C-B netns HOSTC

    # NERT configuration
    run ip netns exec HOSTA ip addr add 172.27.1.1/24 dev veth-A-B
    run ip netns exec HOSTB ip addr add 172.27.1.2/24 dev veth-B-A
    run ip netns exec HOSTB ip addr add 172.27.2.2/24 dev veth-B-C
    run ip netns exec HOSTC ip addr add 172.27.2.1/24 dev veth-C-B

    run ip netns exec HOSTA ip link set veth-A-B up
    run ip netns exec HOSTB ip link set veth-B-A up
    run ip netns exec HOSTB ip link set veth-B-C up
    run ip netns exec HOSTC ip link set veth-C-B up

    run ip netns exec HOSTB sysctl net.ipv4.ip_forward=1

    run ip netns exec HOSTA ip route add 172.27.2.0/24 via 172.27.1.2 dev veth-A-B
    run ip netns exec HOSTC ip route add 172.27.1.0/24 via 172.27.2.2 dev veth-C-B
}

destroy_network () {
    run ip netns del HOSTA
    run ip netns del HOSTB
    run ip netns del HOSTC
}

stop () {
    destroy_network
}

trap stop 0 1 2 3 13 14 15

# exec functions
create_network

status=0; $SHELL || status=$?
exit $status
