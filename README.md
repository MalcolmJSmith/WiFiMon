WiFi Monitor
========

Introduction
--------------

Uses Windows Native Wifi to monitor  signal strength of the connected wireless network and makes this available via the web.

Run wifimon where you would like to maximise the wifi signal. Access its webserver from another device located at the base station. You can then monitor the effect of repositioning the base station and its antennae.

Requirements
----------------

Windows XP or later

Usage
-------

wifimon [port]

port defaults to 80

Output
--------

Shows the currently connected network, the signal quality as  a percentage and the signal strength in db.
Windows Native WiFI considers the signal quality to be 100% when the signal strength is greater than -50db