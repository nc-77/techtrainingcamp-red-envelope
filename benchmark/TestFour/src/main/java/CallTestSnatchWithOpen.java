import org.apache.log4j.Logger;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.SocketTimeoutException;
import java.net.URL;
import java.util.concurrent.CountDownLatch;

public class CallTestSnatchWithOpen implements Runnable{
    private static Logger log = Logger.getLogger(CallTestSnatchWithOpen.class);
    public static int successRequest = 0;
    public static int failRequest = 0;
    public static int timeOutRequest = 0;
    private String url = "http://124.238.238.169:80/v0/snatch";
    private String id = "";
    public static int cur_size_increase = 0;
    public static long costTime = 0;
    public static String envelope_id = "";

    private CountDownLatch begin;
    private CountDownLatch end;

    public static int open_size_increase = 0;
    public static long openCostTime = 0;
    CallTestSnatchWithOpen(String i, CountDownLatch begin, CountDownLatch end){
        this.id = i;
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

    private  static synchronized void incrementCurSizeCount(){
        cur_size_increase++;
    }

    private static synchronized void incrementCostTime(long cost){
        costTime += cost;
    }

    private static synchronized void incrementOpenSizeCount(){ open_size_increase++; }

    private static synchronized void incrementOpenCostTime(long cost){ openCostTime += cost; }

    @Override
    public void run() {
        HttpURLConnection httpURLConnection = null;
        try{
            begin.await();
            try {
                //begin.await();
                long startTime = System.currentTimeMillis();
                // 1. 获取访问地址URL
                URL url = new URL("http://124.238.238.169:80/v0/snatch");
                // 2. 创建HttpURLConnection对象
                httpURLConnection = (HttpURLConnection) url.openConnection();
                /* 3. 设置请求参数等 */
                // 请求方式  默认 GET
                httpURLConnection.setRequestMethod("POST");
                // 超时时间
                httpURLConnection.setConnectTimeout(1000);
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
                String params = "uid="+id;
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
                if(code ==200 && msg.charAt(8)=='0'){
                    incrementCurSizeCount();
                    int index = msg.indexOf("nveloped_id");
                    envelope_id = msg.substring(index+14,index+34);
                    System.out.print("open information:envelope_id is:"+envelope_id+" uid is:"+id+"  ");
                }
                long endTime = System.currentTimeMillis();
                incrementCostTime(endTime-startTime);
                System.out.println("snatch information:id is:"+id+"msg:"+msg);

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
            }

            if(envelope_id!=""){
                try {
                    //begin.await();
                    long startTime = System.currentTimeMillis();
                    // 1. 获取访问地址URL
                    URL url = new URL("http://124.238.238.169:80/v0/open");
                    // 2. 创建HttpURLConnection对象
                    httpURLConnection = (HttpURLConnection) url.openConnection();
                    /* 3. 设置请求参数等 */
                    // 请求方式  默认 GET
                    httpURLConnection.setRequestMethod("POST");
                    // 超时时间
                    httpURLConnection.setConnectTimeout(1000);
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
                        incrementOpenSizeCount();
                    }
                    long endTime = System.currentTimeMillis();
                    incrementOpenCostTime(endTime-startTime);
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
                }
            }

        }catch (Exception e){
            log.error(e.getMessage(),e);
        }finally{
            end.countDown();
        }


    }
}
