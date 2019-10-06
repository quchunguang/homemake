#define PI 3.1415926535897932384626433832795    //定义常数π
 
void setup() {
 
  pinMode(4, ANALOG);     //配置DAC输出
  analogReference(INTERNAL4V096);   //内部基准源4.096V
}
 
void loop() 
{
 for(float i=0;i<=2;i=i+0.01)    //起始点为0,终止为2π，采样率为0.01
 {
  float rad=PI*i;    
  float Sin=sin(rad);
  long intSin=Sin*300;    //将数据放大300倍，取整数
  byte val=map(intSin,-300,300,0,255);   //映射至8位DAC精度
  analogWrite(4, val);   //DAC输出
  }
}
