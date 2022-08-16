package org.vite.portal.orchestrator.beans;

import com.alibaba.fastjson.annotation.JSONField;

import java.util.List;

public class NodeInfo {
    private String sid;
    private OSInfo osInfo;
    private String ip;
    private ProcessInfo processInfo;
    private RuntimeInfo runtimeInfo;
    private long round;
    private boolean isok;

    public long getRound() {
        return round;
    }

    public String getIp() {
        return ip;
    }

    public void setIp(String ip) {
        this.ip = ip;
    }

    public void setRound(long round) {
        this.round = round;
    }

    public boolean isIsok() {
        return isok;
    }

    public void setIsok(boolean isok) {
        this.isok = isok;
    }

    public String getSid() {
        return sid;
    }

    public void setSid(String sid) {
        this.sid = sid;
    }

    public OSInfo getOsInfo() {
        return osInfo;
    }

    public void setOsInfo(OSInfo osInfo) {
        this.osInfo = osInfo;
    }

    public ProcessInfo getProcessInfo() {
        return processInfo;
    }

    public void setProcessInfo(ProcessInfo processInfo) {
        this.processInfo = processInfo;
    }

    public RuntimeInfo getRuntimeInfo() {
        return runtimeInfo;
    }

    public void setRuntimeInfo(RuntimeInfo runtimeInfo) {
        this.runtimeInfo = runtimeInfo;
    }

    /*
    Example go-vite response:
        {
            "jsonrpc":"2.0",
            "id":1,
            "result":{
                "cpuNum":8,
                "gorountine":89,
                "kernelVersion":"18.0.0",
                "os":"darwin",
                "platform":"darwin",
                "platformFamily":"",
                "platformVersion":"10.14"
            }
        }
    **/
    public static class OSInfo {
        private int cpuNum;
        private int gorountine;
        private String kernelVersion;
        private String os;
        private String platform;
        private String platformFamily;
        private String platformVersion;
        private String memTotal;
        private String memFree;

        public int getCpuNum() {
            return cpuNum;
        }

        public void setCpuNum(int cpuNum) {
            this.cpuNum = cpuNum;
        }

        public int getGorountine() {
            return gorountine;
        }

        public void setGorountine(int gorountine) {
            this.gorountine = gorountine;
        }

        public String getKernelVersion() {
            return kernelVersion;
        }

        public void setKernelVersion(String kernelVersion) {
            this.kernelVersion = kernelVersion;
        }

        public String getOs() {
            return os;
        }

        public void setOs(String os) {
            this.os = os;
        }

        public String getPlatform() {
            return platform;
        }

        public void setPlatform(String platform) {
            this.platform = platform;
        }

        public String getPlatformFamily() {
            return platformFamily;
        }

        public void setPlatformFamily(String platformFamily) {
            this.platformFamily = platformFamily;
        }

        public String getPlatformVersion() {
            return platformVersion;
        }

        public void setPlatformVersion(String platformVersion) {
            this.platformVersion = platformVersion;
        }

        public String getMemTotal() {
            return memTotal;
        }

        public void setMemTotal(String memTotal) {
            this.memTotal = memTotal;
        }

        public String getMemFree() {
            return memFree;
        }

        public void setMemFree(String memFree) {
            this.memFree = memFree;
        }
    }

    /*
    Example go-vite response:
        {
            "jsonrpc":"2.0",
            "id":1,
            "result":{
                "build_version":"1.0.2",
                "commit_version":"61917a128a1eab984fea8f2e2d7090af4ffb5b44"
            }
        }
    **/
    public static class ProcessInfo {
        @JSONField(name = "build_version")
        private String buildVersion;

        @JSONField(name = "commit_version")
        private String commitVersion;

        private String rewardAddress;

        private String email;

        private String nodeName;

        private int pid;


        public int getPid() {
            return pid;
        }

        public void setPid(int pid) {
            this.pid = pid;
        }

        public String getBuildVersion() {
            return buildVersion;
        }

        public void setBuildVersion(String buildVersion) {
            this.buildVersion = buildVersion;
        }

        public String getCommitVersion() {
            return commitVersion;
        }

        public void setCommitVersion(String commitVersion) {
            this.commitVersion = commitVersion;
        }

        public String getRewardAddress() {
            return rewardAddress;
        }

        public void setRewardAddress(String rewardAddress) {
            this.rewardAddress = rewardAddress;
        }

        public String getEmail() {
            return email;
        }

        public void setEmail(String email) {
            this.email = email;
        }

        public String getNodeName() {
            return nodeName;
        }

        public void setNodeName(String nodeName) {
            this.nodeName = nodeName;
        }
    }

    /*
    Example go-vite response:
        {
            "jsonrpc":"2.0",
            "id":1,
            "result":{
                "accountPendingNum":"0",
                "latestSnapshot":{
                    "Hash":"0f775a23f1ef1364fa4bf419f4161f2be133d3df2f862470cf27bc715d7b6ed3",
                    "Height":1768708,
                    "Time":1543492685000
                },
                "peersNum":18,
                "producer":"vite_0acbb1335822c8df4488f3eea6e9000eabb0f19d8802f57c87",
                "snapshotPendingNum":0,
                "updateTime":1543492685688
            }
        }
    **/
    public static class RuntimeInfo {
        private String signData;
        private long accountPendingNum;
        private long snapshotPendingNum;
        private long updateTime;
        private long peersNum;
        private String producer;
        private List<Long> delayTime; // [0] Current [1] 1 hour [2] 12 hours [3] 24 hours
        private SnapshotBlock latestSnapshot;

        public List<Long> getDelayTime() {
            return delayTime;
        }

        public void setDelayTime(List<Long> delayTime) {
            this.delayTime = delayTime;
        }

        public String getSignData() {
            return signData;
        }

        public void setSignData(String signData) {
            this.signData = signData;
        }

        public long getAccountPendingNum() {
            return accountPendingNum;
        }

        public void setAccountPendingNum(long accountPendingNum) {
            this.accountPendingNum = accountPendingNum;
        }

        public long getSnapshotPendingNum() {
            return snapshotPendingNum;
        }

        public void setSnapshotPendingNum(long snapshotPendingNum) {
            this.snapshotPendingNum = snapshotPendingNum;
        }

        public long getUpdateTime() {
            return updateTime;
        }

        public void setUpdateTime(long updateTime) {
            this.updateTime = updateTime;
        }

        public long getPeersNum() {
            return peersNum;
        }

        public void setPeersNum(long peersNum) {
            this.peersNum = peersNum;
        }

        public String getProducer() {
            return producer;
        }

        public void setProducer(String producer) {
            this.producer = producer;
        }

        public SnapshotBlock getLatestSnapshot() {
            return latestSnapshot;
        }

        public void setLatestSnapshot(SnapshotBlock latestSnapshot) {
            this.latestSnapshot = latestSnapshot;
        }
    }

    public static class SnapshotBlock {
        @JSONField(name = "Hash")
        private String hash;

        @JSONField(name = "Height")
        private long height;

        @JSONField(name = "Time")
        private long time;

        public String getHash() {
            return hash;
        }

        public void setHash(String hash) {
            this.hash = hash;
        }

        public long getHeight() {
            return height;
        }

        public void setHeight(long height) {
            this.height = height;
        }

        public long getTime() {
            return time;
        }

        public void setTime(long time) {
            this.time = time;
        }
    }
}
