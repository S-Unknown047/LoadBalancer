when you restart your system yout virtual ip will be gone 
sudo ip addr add 192.168.1.100/24 dev wlp0s20f3
execute this command to add virtual ip to your system

sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 80 -j DROP
# This stops the kernel from sending RST to your backend servers' responses
sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 49152:65535 -j DROP

sudo iptables -A INPUT -p tcp -d 127.0.0.2 --dport 49152:65535 -j DROP

sudo iptables -A INPUT -p tcp -d 127.0.0.2 --dport 49152:65535 -j DROP
