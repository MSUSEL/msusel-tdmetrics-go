using Commons.Extensions;
using System;
using System.Collections.Generic;

namespace UnitTests.CommonsTests;

public class ExtensionsTests {

    [Test]
    public void BinarySearchValue() {
        List<int> list = [3, 7, 8, 12, 15, 20];
        checkBinarySearch(list, 1, -1, false);
        checkBinarySearch(list, 2, -1, false);
        checkBinarySearch(list, 3, 0, true);
        checkBinarySearch(list, 4, 0, false);
        checkBinarySearch(list, 6, 0, false);
        checkBinarySearch(list, 7, 1, true);
        checkBinarySearch(list, 8, 2, true);
        checkBinarySearch(list, 9, 2, false);
        checkBinarySearch(list, 12, 3, true);
        checkBinarySearch(list, 13, 3, false);
        checkBinarySearch(list, 18, 4, false);
        checkBinarySearch(list, 19, 4, false);
        checkBinarySearch(list, 20, 5, true);
        checkBinarySearch(list, 21, 5, false);
    }

    private static void checkBinarySearch<T>(IList<T> list, T value, int expIndex, bool expFound)
        where T : IComparable<T> {
        (int gotIndex, bool gotFound) = list.BinarySearch(value);
        Assert.Multiple(() => {
            string msg = "BinarySearch(" + value + ") => (" + gotIndex + ", " + gotFound + ") exp (" + expIndex + ", " + expFound + ")";
            Assert.That(gotIndex, Is.EqualTo(expIndex), "Wrong \"index\" from " + msg);
            Assert.That(gotFound, Is.EqualTo(expFound), "Wrong \"found\" from " + msg);
        });
    }
}
