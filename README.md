### 使用方式   
在config.json里填好相应的信息。  

targetGroupName是想要转发的群名  
 
timeInterval是每次获取帖子的间隔，单位是分钟  

telephone是登软工集市的手机号  
email是登软工集市的邮箱  

password是登软工集市的密码，但是是加密过的，我没有用go复现成功，只能手动复制了。  

密码获取方法：打开开发者工具，输入你的密码和随便什么邮箱，最好是错的邮箱，密码对就行。输入完后登录，查看开发者工具中的网络的login中的载荷，里面有个password,把它的值复制下来就可以了。  

str是用于格式化输出的东西（好像是叫这个？）从左往右依次是用户、标题、标签、pc端链接、移动端链接