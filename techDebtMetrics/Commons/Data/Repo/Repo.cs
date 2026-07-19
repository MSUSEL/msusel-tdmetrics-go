using System;

namespace Commons.Data.Repo;

public class Repo {  
    static public readonly string RepoDir;
    static public readonly string TestDataDir;
    static public readonly string MetricsZip;

    static Repo() {
        const string repoName = "msusel-tdmetrics-go";
        string curDir = Environment.CurrentDirectory;
        int index = curDir.LastIndexOf(repoName);
        if (index == -1) throw new Exception("Failed to find root directory of the repo from " + curDir);
        index += repoName.Length;
        RepoDir = curDir[0..index];
        
        TestDataDir = RepoDir + "/testData";
        MetricsZip  = RepoDir + "/tdd/per_project.zip";
    }

    static public string AbstractedJava(JavaTarget target) =>
        RepoDir + "/tdd/abstractions/" + target.ProjectKey + ".json";

    static public string TestPath(string sourceLang, int testNum) =>
        string.Format("{0}/{1}/test{2:D4}", TestDataDir, sourceLang, testNum);
}
