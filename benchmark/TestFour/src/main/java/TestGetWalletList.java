import org.apache.log4j.Logger;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class TestGetWalletList {

    private static Logger log = Logger.getLogger(TestGetWalletList.class);

    public static void main(String[] args) throws Exception{

        log.debug("This is debug message.");
        // 记录info级别的信息
        log.info("This is info message.");
        // 记录warn级别的信息
        log.info("This is warn message.");
        // 记录error级别的信息
        log.error("This is error message.");

        CountDownLatch begin = new CountDownLatch(1);

        int allRequestSize = 300000;
        log.info("all request size is "+allRequestSize);
        ExecutorService exec = Executors.newFixedThreadPool(1000);
        CountDownLatch end = new CountDownLatch(allRequestSize);
        for(int i=0;i<allRequestSize;i++){
            exec.execute(new CallTestGetWalletList(String.valueOf(i),begin,end));
        }

        long startTime = System.currentTimeMillis();
        begin.countDown();
        try{
            end.await();
        }catch (InterruptedException e){
            e.printStackTrace();
        }finally {
            log.info("all url requests is done!");
            log.info("the success size :"+CallTestGetWalletList.successRequest);
            log.info("the fail size: " + CallTestGetWalletList.failRequest);
            log.info("the timeout size: " + CallTestGetWalletList.timeOutRequest);
            double successRate = (double)CallTestGetWalletList.successRequest / allRequestSize;
            log.info("the success rate is: " + successRate*100+"%");
            long endTime = System.currentTimeMillis();
            long costTime = endTime - startTime;
            log.info("the total time cost is: " + costTime + " ms");
            log.info("qps is:" + allRequestSize/(costTime/1000));
            log.info("sum request time cost is:"+CallTestGetWalletList.costTime+" ms");
        }
        exec.shutdown();
        log.info("main method end");
    }
}
