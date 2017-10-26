create database resistor;
GRANT USAGE ON `resistor`.* to 'resistoruser'@'localhost' identified by 'resistorpass';
GRANT ALL PRIVILEGES ON `resistor`.* to 'resistoruser'@'localhost' with grant option;
flush privileges;
