using Commons.Data.Locations;
using Commons.Data.Yaml;
using System;

namespace UnitTests.CommonsTests;

public class LocationsTests {

    [Test]
    public void LocationRoundTrip() {
        Location loc1 = new(4, "./bob.txt");
        Location loc2 = new(30, "./tommy.txt");
        Location loc3 = new(1, "./jill.txt");
        Location loc4 = new(12, "./bob.txt");
        Location loc5 = new(6, "./jill.txt");

        Factory factory = new();
        factory.Add(loc1);
        factory.Add(loc2);
        factory.Add(loc3);
        factory.Add(loc4);
        factory.Add(loc5);

        Writer writer = factory.Build();
        Commons.Data.Yaml.Object obj = new();
        obj.Add("locs", writer.Write());
        obj.Add("loc1", writer, loc1);
        obj.Add("loc2", writer, loc2);
        obj.Add("loc3", writer, loc3);
        obj.Add("loc4", writer, loc4);
        obj.Add("loc5", writer, loc5);

        string yaml = obj.ToString();
        string[] gotYaml = yaml.Trim().Split("\n");
        string[] expYaml = [
            "locs:",
            "  0: ./bob.txt",
            "  12: ./jill.txt",
            "  18: ./tommy.txt",
            "loc1: 3",
            "loc2: 47",
            "loc3: 12",
            "loc4: 11",
            "loc5: 17",
            "..."];
        if (gotYaml != expYaml) {
            Console.WriteLine(gotYaml);
            Assert.That(gotYaml, Is.EqualTo(expYaml));
        }

        Commons.Data.Yaml.Object root = Node.Parse(yaml).AsObject();
        Reader reader = Reader.Read(root.ReadNode("locs"));
        Location got1 = root.ReadLocation(reader, "loc1");
        Location got2 = root.ReadLocation(reader, "loc2");
        Location got3 = root.ReadLocation(reader, "loc3");
        Location got4 = root.ReadLocation(reader, "loc4");
        Location got5 = root.ReadLocation(reader, "loc5");

        Assert.Multiple(() => {
            Assert.That(got1.ToString(), Is.EqualTo("./bob.txt:4"));
            Assert.That(got2.ToString(), Is.EqualTo("./tommy.txt:30"));
            Assert.That(got3.ToString(), Is.EqualTo("./jill.txt:1"));
            Assert.That(got4.ToString(), Is.EqualTo("./bob.txt:12"));
            Assert.That(got5.ToString(), Is.EqualTo("./jill.txt:6"));

            Assert.That(reader[-10].ToString(), Is.EqualTo("<unknown>:0"));
            Assert.That(reader[-1].ToString(), Is.EqualTo("<unknown>:0"));
            Assert.That(reader[0].ToString(), Is.EqualTo("./bob.txt:1"));
            Assert.That(reader[11].ToString(), Is.EqualTo("./bob.txt:12"));
            Assert.That(reader[12].ToString(), Is.EqualTo("./jill.txt:1"));
            Assert.That(reader[17].ToString(), Is.EqualTo("./jill.txt:6"));
            Assert.That(reader[18].ToString(), Is.EqualTo("./tommy.txt:1"));
            Assert.That(reader[300].ToString(), Is.EqualTo("./tommy.txt:283"));
        });
    }
}
