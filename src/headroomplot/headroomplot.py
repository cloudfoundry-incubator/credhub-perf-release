import sys
import os.path
import re
import warnings
from io import StringIO

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import statsmodels.api as sm
from patsy import dmatrices

matplotlib.style.use('ggplot')

matplotlib.rcParams['figure.figsize'] = 9, 6
matplotlib.rcParams['legend.loc'] = 'best'
matplotlib.rcParams['figure.dpi'] = 120

# We'll need these packages for plotting fit lines
warnings.filterwarnings('ignore')
performanceResultsFile = sys.argv[1]
assert os.path.isfile(performanceResultsFile), 'Missing performance results file'

compareDatasets = False

if compareDatasets:
    assert os.path.isfile('old_perfResults.csv'), 'Missing old performance results file "old_perfResults.csv"'


def readThroughputData(filename):
    with open(filename) as f:
        data = f.read()
    # Get locations of start-time,response-time headers in file
    header_idxs = [m.start() for m in re.finditer('start-time,response-time', data)]
    header_idxs.append(len(data))
    prev = header_idxs[0]

    df = pd.DataFrame()
    # Read each section delimited by the csv headers
    for cur in header_idxs[1:]:
        dfSection = pd.read_csv(StringIO(unicode(data[prev:cur])), parse_dates=['start-time'])
        trimmedSection = trimEdges(dfSection)

        if len(trimmedSection) == 0:
            print "There is not enough data to build headroom plot. Please increase the number of requests."
            exit(1)

        df = df.append(trimmedSection)
        prev = cur
    # Reset the index because it is a Frankenstein of smaller indexes
    df = df.reset_index().drop('index', axis=1)
    return df


def trimEdges(data):
    indexes = data.set_index('start-time').resample('1S').aggregate(lambda x: 1).index
    testStartTime = indexes[1]
    testEndTime = indexes[-2]
    return data[(data['start-time'] >= testStartTime) & (data['start-time'] <= testEndTime)]


def processThroughputData(data):
    buckets = data.set_index('start-time')['response-time'].resample('1S')
    throughputDataSet = buckets.aggregate({"throughput": lambda x: 0 if x.count() == 0 else x.count()})

    throughputDataSet = throughputDataSet.reset_index()
    throughputDataSet = throughputDataSet.fillna(method='ffill')
    return buckets, throughputDataSet


goData = readThroughputData(performanceResultsFile)

throughputBuckets, throughputData = processThroughputData(goData)

if compareDatasets:
    oldGoData = readThroughputData('old_perfResults.csv')
    oldThroughputBuckets, oldThroughputData = processThroughputData(oldGoData)

goData['throughput'] = throughputBuckets.transform(len).reset_index()['response-time']
goData.columns = ['start-time', 'latency', 'throughput']

if compareDatasets:
    oldGoData['throughput'] = oldThroughputBuckets.transform(len).reset_index()['response-time']
    oldGoData.columns = ['start-time', 'latency', 'throughput']


def generateFitLine(data):
    y, x = dmatrices('latency ~ throughput', data=data, return_type='dataframe')
    fit = sm.GLM(y, x, family=sm.families.InverseGaussian(sm.families.links.inverse_squared)).fit()
    maxThroughput = data['throughput'].max()
    minThroughtput = data['throughput'].min()
    domain = np.arange(minThroughtput, maxThroughput)
    predictionInputs = np.ones((len(domain), 2))
    predictionInputs[:, 1] = domain
    fitLine = fit.predict(predictionInputs)
    return domain, fitLine, round(maxThroughput)


domain, goFitLine, xLimit = generateFitLine(goData)

if compareDatasets:
    oldDomain, oldGoFitLine, oldXLimit = generateFitLine(oldGoData)

fig, ax = plt.subplots()

# Change the value of `c` to change the color. http://matplotlib.org/api/colors_api.html
ax = goData.plot(ax=ax, kind='scatter', x='throughput', y='latency', c='b', marker='.', alpha=0.2)
ax.plot(domain, goFitLine, c='b', lw=2)  # Plot the fit line

if compareDatasets:
    ax = oldGoData.plot(ax=ax, kind='scatter', x='throughput', y='latency', c='r', marker='.', alpha=0.2)
    ax.plot(oldDomain, oldGoFitLine, c='r', lw=2)  # Plot the fit line
    ax.legend(['after', 'before'])

# To update x & y axis range change the parameters in function set_(x/y)lim(lower_limit, uppper_limit)

ax.set_ylim(0, 10)
ax.set_xlim(0, xLimit)
ax.autoscale_view(True, True, True)
plt.xlabel('Throughput (requests/sec)')
plt.ylabel('Latency (sec)')
plt.title('Headroom plot', y=1.05)
plt.plot()

filenameForPlot = performanceResultsFile[:-4] + "Plot.png"
plt.savefig(filenameForPlot)
print ("saving graph to " + filenameForPlot)
