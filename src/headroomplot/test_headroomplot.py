import os

from unittest import TestCase
from headroomplot import readThroughputData


class TestReadThroughputData(TestCase):
    THIS_DIR = os.path.dirname(os.path.abspath(__file__))

    def test_fails_when_file_is_not_found(self):
        nonexistent_file = "/tmp/does_not_exit"
        with self.assertRaises(IOError) as context:
            readThroughputData(nonexistent_file)

    def test_return_empty_when_file_does_not_contain_correct_header_names(self):
        wrong_headers_file = os.path.join(self.THIS_DIR, 'test_data/wrong_headers.csv')
        data = readThroughputData(wrong_headers_file)
        self.assertEquals(len(data), 0)

    def test_return_data_frame_with_the_last_second_trimmed_out(self):
        working_file = os.path.join(self.THIS_DIR, 'test_data/working.csv')
        data = readThroughputData(working_file)
        self.assertEquals(len(data), 6)

    def test_return_empty_when_test_date_is_only_one_second(self):
        a_second_only_file = os.path.join(self.THIS_DIR, 'test_data/a_second_only.csv')
        data = readThroughputData(a_second_only_file)
        self.assertEquals(len(data), 0)
