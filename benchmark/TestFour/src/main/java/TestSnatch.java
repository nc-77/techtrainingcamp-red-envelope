
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import org.apache.log4j.Logger;


public class TestSnatch {

    private static Logger log = Logger.getLogger(TestSnatch.class);

    public static void main(String[] args) throws Exception{

        log.debug("This is debug message.");
        // 记录info级别的信息
        log.info("This is info message.");
        // 记录warn级别的信息
        log.info("This is warn message.");
        // 记录error级别的信息
        log.error("This is error message.");

        CountDownLatch begin = new CountDownLatch(1);

        int allRequestSize = 100000;
        log.info("all request size is "+allRequestSize);
        ExecutorService exec = Executors.newFixedThreadPool(1000);
        CountDownLatch end = new CountDownLatch(allRequestSize);
        for(int j=0;j<15;j++){
            for(int i=0;i<allRequestSize/15;i++){
                exec.execute(new CallTestSnatch(String.valueOf(i+100800000),begin,end));
            }
        }
        long startTime = System.currentTimeMillis();
        begin.countDown();
        try{
            end.await();
        }catch (InterruptedException e){
            e.printStackTrace();
        }finally {
            log.info("all url requests is done!");
            log.info("the success size :"+CallTestSnatch.successRequest);
            log.info("the fail size: " + CallTestSnatch.failRequest);
            log.info("the timeout size: " + CallTestSnatch.timeOutRequest);
            double successRate = (double)CallTestSnatch.successRequest / allRequestSize;
            log.info("the success rate is: " + successRate*100+"%");
            long endTime = System.currentTimeMillis();
            long costTime = endTime - startTime;
            log.info("the total time cost is: " + costTime + " ms");
            log.info("qps is:" + allRequestSize/(costTime/1000));
            log.info("cur_size_increase is: "+CallTestSnatch.cur_size_increase);
            log.info("sum request time cost is:"+CallTestSnatch.costTime+" ms");
        }
        exec.shutdown();
        log.info("main method end");
    }
}
