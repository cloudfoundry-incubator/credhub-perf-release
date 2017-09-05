import os

from unittest import TestCase
from headroomplot import read_throughput_data
from headroomplot import PerformanceRunIterator
from headroomplot import PerfData


class TestPerfData(TestCase):
    THIS_DIR = os.path.dirname(os.path.abspath(__file__))

    def setUp(self):
        self.working_file = os.path.join(self.THIS_DIR, 'test_data/working.csv')

    def test_data_return_file_content(self):
        perf_data = PerfData(self.working_file)
        self.assertEquals("start-time,response-time\n2017-08-25T21:07:33.0080871Z,0.0775\n2017-08-25T21:07:33.0855854Z,0.0272\n2017-08-25T21:07:34.001946224Z,0.0249\n2017-08-25T21:07:34.02683099Z,0.0378\nstart-time,response-time\n2017-08-25T21:07:35.114091565Z,0.029\n2017-08-25T21:07:36.144091565Z,0.029\n2017-08-25T21:07:37.144091565Z,0.029", perf_data.data())

    def test_header_returns_perf_header(self):
        perf_data = PerfData(self.working_file)
        self.assertEquals("start-time,response-time", perf_data.headers())

    def test_datetime_headers_returns_headers_that_hold_datetime(self):
        perf_data = PerfData(self.working_file)
        self.assertEquals("start-time", perf_data.datetime_headers())


class TestPerformanceRunIterator(TestCase):

    def setUp(self):
        self._header = "start-time,response-time"

    def test_works_with_one_run(self):
        data = self._header +"\n2017-08-25T21:07:33.0080871Z,0.0775"

        for run in PerformanceRunIterator(data, self._header):
            self.assertEquals(data, run)

    def test_works_with_more_than_one_run(self):

        run_item = self._header +"\n2017-08-25T21:07:33.0080871Z,0.0775\n"
        run_item2 = self._header + "\n2017-08-25T21:07:33.0080871Z,0.0779\n"
        data = run_item + run_item2
        run_items = [run_item, run_item2]

        index = 0
        for run in PerformanceRunIterator(data, self._header):
            self.assertEquals(run_items[index], run)
            index = index + 1

        self.assertEquals(2, index)


class TestReadThroughputData(TestCase):
    THIS_DIR = os.path.dirname(os.path.abspath(__file__))

    def test_fails_when_file_is_not_found(self):
        nonexistent_file = "/tmp/does_not_exit"
        with self.assertRaises(IOError) as context:
            read_throughput_data(nonexistent_file)

    def test_return_empty_when_file_does_not_contain_correct_header_names(self):
        wrong_headers_file = os.path.join(self.THIS_DIR, 'test_data/wrong_headers.csv')
        data = read_throughput_data(wrong_headers_file)
        self.assertEquals(len(data), 0)

    def test_return_data_frame_with_the_last_second_trimmed_out(self):
        working_file = os.path.join(self.THIS_DIR, 'test_data/working.csv')
        data = read_throughput_data(working_file)
        self.assertEquals(len(data), 4)

    def test_return_empty_when_test_date_is_only_one_second(self):
        a_second_only_file = os.path.join(self.THIS_DIR, 'test_data/a_second_only.csv')
        data = read_throughput_data(a_second_only_file)
        self.assertEquals(len(data), 0)
