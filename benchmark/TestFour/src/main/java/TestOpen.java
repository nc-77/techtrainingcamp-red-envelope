
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;


import org.apache.log4j.Logger;

public class TestOpen {

    private static Logger log = Logger.getLogger(TestOpen.class);

    public static void main(String[] args) throws Exception{

        log.debug("This is debug message.");
        // 记录info级别的信息
        log.info("This is info message.");
        // 记录warn级别的信息
        log.info("This is warn message.");
        // 记录error级别的信息
        log.error("This is error message.");

        CountDownLatch begin = new CountDownLatch(1);

        int allRequestSize = 2000;
        String envelope_id = "c67479bbu3iato50glpg";
        log.info("all request size is "+allRequestSize);
        ExecutorService exec = Executors.newFixedThreadPool(100);
        CountDownLatch end = new CountDownLatch(allRequestSize);
        for(int i=1;i<1+allRequestSize;i++){
            exec.execute(new CallTestOpen(String.valueOf(150000129),envelope_id,begin,end));
        }

        long startTime = System.currentTimeMillis();
        begin.countDown();
        try{
            end.await();
        }catch (InterruptedException e){
            e.printStackTrace();
        }finally {
            log.info("all url requests is done!");
            log.info("the success size :"+CallTestOpen.successRequest);
            log.info("the fail size: " + CallTestOpen.failRequest);
            log.info("the timeout size: " + CallTestOpen.timeOutRequest);
            double successRate = (double)CallTestOpen.successRequest / allRequestSize;
            log.info("the success rate is: " + successRate*100+"%");
            long endTime = System.currentTimeMillis();
            long costTime = endTime - startTime;
            log.info("the total time cost is: " + costTime + " ms");
            log.info("average request time cost is: " + costTime / allRequestSize + " ms");
            log.info("qps is:" + allRequestSize/(costTime/1000));
            log.info("open_size_increase is: "+CallTestOpen.open_size_increase);
        }
        exec.shutdown();
        log.info("main method end");
    }
}
