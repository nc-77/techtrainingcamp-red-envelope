import org.apache.log4j.Logger;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.SocketTimeoutException;
import java.net.URL;
import java.util.concurrent.CountDownLatch;

public class CallTestOpen implements Runnable{
    private static Logger log = Logger.getLogger(CallTestOpen.class);
    public static int successRequest = 0;
    public static int failRequest = 0;
    public static int timeOutRequest = 0;
    private String url = "http://124.238.238.169:80/v0/snatch";
    private String id = "";
    private String envelope_id = "";
    public static int open_size_increase = 0;    //成功打开的红包个数

    private CountDownLatch begin;
    private CountDownLatch end;
    CallTestOpen(String i,String envelope_id, CountDownLatch begin, CountDownLatch end){
        this.id = i;
        this.envelope_id = envelope_id;
        this.begin = begin;
        this.end = end;
    }

    private  static synchronized void incrementSuccessCount(){
        successRequest++;
    }

    private  static synchronized void incrementFailCount(){
        failRequest++;
    }

    private static synchronized void incrementTimeOutCount(){
        timeOutRequest++;
    }

    private  static synchronized void incrementcursizeCount(){
        open_size_increase++;
    }

    @Override
    public void run() {
        HttpURLConnection httpURLConnection = null;
        try {
            begin.await();
            // 1. 获取访问地址URL
            URL url = new URL("http://124.238.238.169:80/v0/open");
            // 2. 创建HttpURLConnection对象
            httpURLConnection = (HttpURLConnection) url.openConnection();
            /* 3. 设置请求参数等 */
            // 请求方式  默认 GET
            httpURLConnection.setRequestMethod("POST");
            // 超时时间
            httpURLConnection.setConnectTimeout(300);
            // 设置是否输出
            httpURLConnection.setDoOutput(true);
            // 设置是否读入
            httpURLConnection.setDoInput(true);
            // 设置是否使用缓存
            httpURLConnection.setUseCaches(false);
            // 设置此 HttpURLConnection 实例是否应该自动执行 HTTP 重定向
            httpURLConnection.setInstanceFollowRedirects(true);
            // 设置请求头
            httpURLConnection.addRequestProperty("sysId","sysId");
            // 设置使用标准编码格式编码参数的名-值对
            httpURLConnection.setRequestProperty("Content-Type", "application/x-www-form-urlencoded");
            // 连接
            httpURLConnection.connect();
            /* 4. 处理输入输出 */
            // 写入参数到请求中
            String params = "uid=" + id +"&envelope_id=" + envelope_id;
            OutputStream out = httpURLConnection.getOutputStream();
            out.write(params.getBytes());
            // 简化
            //httpURLConnection.getOutputStream().write(params.getBytes());
            out.flush();
            out.close();
            // 从连接中读取响应信息
            String msg = "";
            int code = httpURLConnection.getResponseCode();
            if (code == 200) {
                incrementSuccessCount();
                BufferedReader reader = new BufferedReader(
                        new InputStreamReader(httpURLConnection.getInputStream()));
                String line;
                while ((line = reader.readLine()) != null) {
                    msg += line + "\n";
                }
                reader.close();
            }else{
                incrementFailCount();
            }
            if(code ==200 && msg.charAt(8)=='0'){    //code=1表示这个红包已经打开过，code=0表示这个红包没有打开过
                incrementcursizeCount();
            }
            System.out.println(msg);

        } catch (SocketTimeoutException e) {
            incrementTimeOutCount();
            log.error(e.getMessage(),e);
        }catch (Exception e){
            log.error(e.getMessage(),e);

        }finally{
            // 5. 断开连接
            if (null != httpURLConnection){
                try {
                    httpURLConnection.disconnect();
                }catch (Exception e){
                    e.printStackTrace();
                }
            }
            end.countDown();
        }
    }
}

