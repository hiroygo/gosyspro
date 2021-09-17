## netns で動作を試す
マルチキャスト通信はデフォルトゲートウェイを設定しなくても動作する。

```
$ sudo ip netns add ns1
$ sudo ip netns add ns2
$ sudo ip link add ns1-veth0 type veth peer name ns2-veth0
$ sudo ip link set ns1-veth0 netns ns1
$ sudo ip link set ns2-veth0 netns ns2
$ sudo ip netns exec ns1 ip link set ns1-veth0 up
$ sudo ip netns exec ns2 ip link set ns2-veth0 up
$ sudo ip netns exec ns1 ip a add 192.0.1.1/24 dev ns1-veth0 
$ sudo ip netns exec ns2 ip a add 192.0.1.2/24 dev ns2-veth0 
$ sudo ip netns exec ns1 route add -net 224.0.0.0/4 dev ns1-veth0
$ sudo ip netns exec ns2 route add -net 224.0.0.0/4 dev ns2-veth0
$ sudo ip netns exec ns1 ./server
$ sudo ip netns exec ns1 ./client
```
