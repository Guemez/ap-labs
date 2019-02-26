Lab 2.1 - ClockWall
===================

To compile
---------------------------------------
```
make
```

To run
---------------------------------------
First open the ports with clock2
For just 1 port
```
TZ=US/Eastern    ./clock2 -port 8010
```
For multiple ports
```
TZ=US/Eastern    ./clock2 -port 8010 &  TZ=Asia/Tokyo    ./clock2 -port 8020 & TZ=Europe/London ./clock2 -port 8030 &
```

Now to listen with clockWall, pass the clock for each port opened with clock2
```
./clockWall NewYork=localhost:8010 Tokyo=localhost:8020 London=localhost:8030
```
