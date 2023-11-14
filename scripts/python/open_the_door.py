import RPi.GPIO as GPIO
import time

# 设置引脚编号模式为BCM模式
GPIO.setmode(GPIO.BCM)

# 定义引脚号
pin17 = 17
pin27 = 27

# 设置引脚为输出模式
GPIO.setup(pin17, GPIO.OUT)
GPIO.setup(pin27, GPIO.OUT)

# 输出高电平到引脚17，低电平到引脚27
GPIO.output(pin17, GPIO.HIGH)
GPIO.output(pin27, GPIO.LOW)

# 步进机给电 3 秒
time.sleep(3)

# 两个GPIO都输出低电平，步进机挺转
GPIO.output(pin17, GPIO.LOW)
GPIO.output(pin27, GPIO.LOW)

# 留 20 秒开门时间
time.sleep(20)


# 输出低电平到引脚17，高电平到引脚27
GPIO.output(pin17, GPIO.LOW)
GPIO.output(pin27, GPIO.HIGH)

# 等待 3 秒
time.sleep(3)

# 清理GPIO资源
GPIO.cleanup()