import org.apache.log4j.Logger;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class TestSnatchWithOpen {
    private static Logger log = Logger.getLogger(TestSnatchWithOpen.class);

    public static void main(String[] args) throws Exception{

        log.debug("This is debug message.");
        // 记录info级别的信息
        log.info("This is info message.");
        // 记录warn级别的信息
        log.info("This is warn message.");
        // 记录error级别的信息
        log.error("This is error message.");

        CountDownLatch begin = new CountDownLatch(1);

        int allRequestSize = 1000;
        log.info("all request size is "+allRequestSize);
        ExecutorService exec = Executors.newFixedThreadPool(10);
        CountDownLatch end = new CountDownLatch(allRequestSize);
        for(int j=0;j<10;j++){
            for(int i=0;i<allRequestSize/10;i++){
                exec.execute(new CallTestSnatchWithOpen(String.valueOf(i+150000990),begin,end));
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
            //log.info("qps is:" + allRequestSize/(costTime/1000));
            log.info("snatch_size_increase is: "+CallTestSnatchWithOpen.cur_size_increase);
            log.info("open_size_increase is: "+ CallTestSnatchWithOpen.open_size_increase);
            log.info("sum snatch request time cost is:"+CallTestSnatchWithOpen.costTime+" ms");
            log.info("sum open request time cost is:" + CallTestSnatchWithOpen.openCostTime+ " ms");
        }
        exec.shutdown();
        log.info("main method end");
    }
}
