using System.Linq;

namespace Commons.Data.Repo;

public readonly record struct JavaTarget(string ProjectKey, string GetLink, string CommitSha, string GroupId) {
    static public readonly JavaTarget Archiva              = new("archiva",               "apache/archiva",               "374fc983abc92df8aa4f8ef30caee94b34312ad2", "org.apache.archiva");
    static public readonly JavaTarget Batik                = new("batik",                 "apache/batik",                 "2bb3a6ea5a6258ff6372e2493b81d7768d6bb494", "org.apache.batik");
    static public readonly JavaTarget Cayenne              = new("cayenne",               "apache/cayenne",               "b9988a83e364b9b470873dff8996dcf401d08dc4", "org.apache.cayenne");
    static public readonly JavaTarget Cocoon               = new("cocoon",                "apache/cocoon",                "a80f73b27592a2794c9133ee03d2e402bf12ecc1", "org.apache.cocoon");
    static public readonly JavaTarget CommonsBcel          = new("commons-bcel",          "apache/commons-bcel",          "6ed18c5bef0f5b93b54783a8e8fb2b9042da26ac", "org.apache.bcel");
    static public readonly JavaTarget CommonsBeanutils     = new("commons-beanutils",     "apache/commons-beanutils",     "c4da598872233b59af41a221bd2bdcefbbca1259", "org.apache.commons.beanutils2");
    static public readonly JavaTarget CommonsCli           = new("commons-cli",           "apache/commons-cli",           "92f1def0bb3c0345295012e36b7150cfd1d7b6ab", "org.apache.commons.cli");
    static public readonly JavaTarget CommonsCodec         = new("commons-codec",         "apache/commons-codec",         "db51a1cb41e9155ca028a73b0637b32a2c37c43a", "org.apache.commons.codec");
    static public readonly JavaTarget CommonsCollections   = new("commons-collections",   "apache/commons-collections",   "f0f364fd9d946483f947011a3557c1e6f2e5d8ee", "org.apache.commons.collections4");
    static public readonly JavaTarget CommonsConfiguration = new("commons-configuration", "apache/commons-configuration", "15b4031ba94a60f20b854e6ce2c7964d77086387", "org.apache.commons.configuration2");
    static public readonly JavaTarget CommonsDaemon        = new("commons-daemon",        "apache/commons-daemon",        "1ffa799cb3ddf5a4a918e59e46cd9868ee766b19", "org.apache.commons.daemon");
    static public readonly JavaTarget CommonsDbcp          = new("commons-dbcp",          "apache/commons-dbcp",          "d8dd39b32bbb04a28ea86eb826c56aa6783f3faf", "org.apache.commons.dbcp2");
    static public readonly JavaTarget CommonsDbutils       = new("commons-dbutils",       "apache/commons-dbutils",       "2f48485a82697d9aed060ba36f6d5beb3a58ed8b", "org.apache.commons.dbutils");
    static public readonly JavaTarget CommonsDigester      = new("commons-digester",      "apache/commons-digester",      "c1d0e563339faec040eb036ae97a7b7bf07ba865", "org.apache.commons.digester3");
    static public readonly JavaTarget CommonsExec          = new("commons-exec",          "apache/commons-exec",          "2da60ab3eefaaa2f8a434ded1eebe1ce17efd34a", "org.apache.commons.exec");
    static public readonly JavaTarget CommonsFileUpload    = new("commons-fileupload",    "apache/commons-fileupload",    "cae90facebc54803232a0593003914ca77193a73", "org.apache.commons.fileupload");
    static public readonly JavaTarget CommonsIo            = new("commons-io",            "apache/commons-io",            "65c4a9c0ec651dd99f28b9fae40378728d071985", "org.apache.commons.io");
    static public readonly JavaTarget CommonsJelly         = new("commons-jelly",         "apache/commons-jelly",         "48c008cc2328402e0976295625b32c5197ba2324", "org.apache.commons.jelly");
    static public readonly JavaTarget CommonsJexl          = new("commons-jexl",          "apache/commons-jexl",          "d3e702149a3db297d6db2c0b7671807f5c7b98fc", "org.apache.commons.jexl3");
    static public readonly JavaTarget CommonsJxpath        = new("commons-jxpath",        "apache/commons-jxpath",        "eff47ab8ca52fdbc91d1313cc224324465dd043e", "org.apache.commons.jxpath");
    static public readonly JavaTarget CommonsNet           = new("commons-net",           "apache/commons-net",           "fb7aae4c64f7d2bf6dced00c49c3ffc428b2d572", "org.apache.commons.net");
    static public readonly JavaTarget CommonsOgnl          = new("commons-ognl",          "apache/commons-ognl",          "6ec1a1a4588b82c0972ca2ff35b85d9b50cc4604", "org.apache.commons.ognl");
    static public readonly JavaTarget CommonsValidator     = new("commons-validator",     "apache/commons-validator",     "a3771313c9f1833abf32c7c294ad1de4810e532d", "org.apache.commons.validator");
    static public readonly JavaTarget CommonsVfs           = new("commons-vfs",           "apache/commons-vfs",           "d72192f18bfaed730b4f37a2f94853e1503ffd74", "org.apache.commons.vfs2");
    static public readonly JavaTarget Felix                = new("felix",                 "apache/felix",                 "bdb6cb5cac0d81e9cd3fda666065e0e577eb9c41", "org.apache.felix");
    static public readonly JavaTarget Hive                 = new("hive",                  "apache/hive",                  "a4d91eaf2925239aa29342f7e5b0f8680c842390", "org.apache.hadoop.hive");
    static public readonly JavaTarget HttpComponentsClient = new("httpcomponents-client", "apache/httpcomponents-client", "8a1b96bfa75382c0b94d70f6914fbb9bfeb0451e", "org.apache.hc.client5");
    static public readonly JavaTarget HttpComponentsCore   = new("httpcomponents-core",   "apache/httpcomponents-core",   "3a677d47cb872b6ede20b28e93d3206f08b349ac", "org.apache.hc.core5");
    static public readonly JavaTarget Santuario            = new("santuario",             "apache/santuario-java",        "be4e2331f77adb1e479406ebf973e516bbf5e32b", "org.apache.jcp");
    static public readonly JavaTarget Thrift               = new("thrift",                "apache/thrift",                "a2123693838410c1e78170419e9bb91cb01151b4", "org.apache.thrift");
    static public readonly JavaTarget Zookeeper            = new("zookeeper",             "apache/zookeeper",             "eac693cc76a34f96b9116ef33d1e92af7129416d", "org.apache.zookeeper");

    static public readonly JavaTarget[] Targets = [
        Archiva,
        Batik,
        Cayenne,
        Cocoon,
        CommonsBcel,
        CommonsBeanutils,
        CommonsCli,
        CommonsCodec,
        CommonsCollections,
        CommonsConfiguration,
        CommonsDaemon,
        CommonsDbcp,
        CommonsDbutils,
        CommonsDigester,
        CommonsExec,
        CommonsFileUpload,
        CommonsIo,
        CommonsJelly,
        CommonsJexl,
        CommonsJxpath,
        CommonsNet,
        CommonsOgnl,
        CommonsValidator,
        CommonsVfs,
        Felix,
        Hive,
        HttpComponentsClient,
        HttpComponentsCore,
        Santuario,
        Thrift,
        Zookeeper,
    ];

    static public JavaTarget? Find(string project_key) =>
        (from t in Targets where t.ProjectKey == project_key select t).First();
}
