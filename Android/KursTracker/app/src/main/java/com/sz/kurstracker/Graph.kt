package com.sz.kurstracker

import android.R
import android.graphics.Color
import android.graphics.DashPathEffect
import android.graphics.Paint
import com.androidplot.ui.DynamicTableModel
import com.androidplot.ui.TableOrder
import com.androidplot.xy.*
import java.text.*
import java.util.*
import kotlin.collections.HashMap


object Graph {
    private val FORMATTERS: Map<String, LineAndPointFormatter>

    init {
        val dashPaint = Paint()
        dashPaint.color = Color.RED
        dashPaint.style = Paint.Style.STROKE
        dashPaint.strokeWidth = 3f
        dashPaint.pathEffect = DashPathEffect(floatArrayOf(10f, 5f), 0.0f)
        val dashPaint2 = Paint()
        dashPaint2.color = Color.RED
        dashPaint2.style = Paint.Style.STROKE
        dashPaint2.strokeWidth = 3f
        dashPaint2.pathEffect = DashPathEffect(floatArrayOf(10f, 5f), 0.0f)
        FORMATTERS = mapOf(
                Pair("MBB", LineAndPointFormatter(Color.RED, Color.GREEN, Color.TRANSPARENT, null)),
                Pair("MBS", LineAndPointFormatter(Color.YELLOW, Color.BLUE, Color.TRANSPARENT, null)),
                Pair("PBB", LineAndPointFormatter(Color.CYAN, Color.MAGENTA, Color.TRANSPARENT, null)),
                Pair("PBS", LineAndPointFormatter(Color.BLACK, Color.valueOf(128f, 0f, 128f).toArgb(), Color.TRANSPARENT, null)),
                Pair("ALB", LineAndPointFormatter(Color.GREEN, Color.CYAN, Color.TRANSPARENT, null)),
                Pair("ALS", LineAndPointFormatter(Color.BLUE, Color.YELLOW, Color.TRANSPARENT, null)),
                Pair("MOB", LineAndPointFormatter(Color.MAGENTA, Color.CYAN, Color.TRANSPARENT, null)),
                Pair("MOS", LineAndPointFormatter(Color.valueOf(128f, 0f, 128f).toArgb(), Color.BLACK, Color.TRANSPARENT, null))
        )
        FORMATTERS["MBB"]!!.linePaint = dashPaint;
        FORMATTERS["MBS"]!!.linePaint = dashPaint2;
    }

    private val FORMATTER = SimpleDateFormat("HH:mm", Locale.US)
    private val FORMATTER2 = SimpleDateFormat("dd.MM", Locale.US)
    private val NUMBER_FORMATTER = DecimalFormat("#.00")

    private val LABEL_FORMATTER = object : Format() {
        override fun format(obj: Any, toAppendTo: StringBuffer, pos: FieldPosition): StringBuffer {
            val idx = (obj as Number).toInt()
            if (mDays <= 2) {
                return toAppendTo.append(FORMATTER.format(mIdxToDate[idx]!!))
            }
            return toAppendTo.append(FORMATTER2.format(mIdxToDate[idx]!!))
        }

        override fun parseObject(source: String, pos: ParsePosition): Any? {
            return null
        }
    }

    private class IdxComparator : Comparator<IGraphData> {
        override fun compare(o1: IGraphData, o2: IGraphData): Int {
            return o1.idx - o2.idx
        }
    }

    private var mIdxToDate: Map<Int, Date> = mapOf()
    private var mDays = 0

    private class GraphDataSeries constructor(private val mData: List<IGraphData>, private val mDataName: String) : OrderedXYSeries {

        override fun size(): Int {
            return mData.size
        }

        override fun getX(index: Int): Number {
            return mData[index].idx
        }

        override fun getY(index: Int): Number? {
            return mData[index].getData(mDataName)
        }

        override fun getTitle(): String {
            return mData[0].seriesColumnData
        }

        override fun getXOrder(): OrderedXYSeries.XOrder {
            return OrderedXYSeries.XOrder.ASCENDING
        }
    }

    private data class SeriesInfo(val mSeries: List<XYSeries>, val mLowerBoundary: Double, val mUpperBoundary: Double)

    private fun buildSeries(data: List<IGraphData>, dataName: String): SeriesInfo {
        val series = HashMap<String, MutableList<IGraphData>>()
        var lowerBoundary = java.lang.Double.MAX_VALUE
        var upperBoundary = java.lang.Double.MIN_VALUE

        for (dataItem in data) {
            var values: MutableList<IGraphData>? = series[dataItem.seriesColumnData]
            if (values == null) {
                values = ArrayList()
                series[dataItem.seriesColumnData] = values
            }

            val dataValue = dataItem.getData(dataName)
            if (dataValue != null) {
                values.add(dataItem)
                if (dataValue < lowerBoundary) {
                    lowerBoundary = dataValue
                }
                if (dataValue > upperBoundary) {
                    upperBoundary = dataValue
                }
            }
        }

        val result = ArrayList<XYSeries>()

        for ((_, value) in series) {
            if (value.size > 1) {
                result.add(GraphDataSeries(value, dataName))
            }
        }

        return SeriesInfo(result, lowerBoundary, upperBoundary)
    }

    fun buildGraph(plot: XYPlot, data: List<IGraphData>, dataName: String) {
        plot.clear()
        if (data.size > 1) {
            val comparator = IdxComparator()
            Collections.sort(data, comparator)
            val seriesInfo = buildSeries(data, dataName)
            if (seriesInfo.mSeries.isNotEmpty()) {
                for (seriesItem in seriesInfo.mSeries) {
                    plot.addSeries(seriesItem, FORMATTERS[seriesItem.title])
                }
                plot.setRangeBoundaries(seriesInfo.mLowerBoundary, seriesInfo.mUpperBoundary, BoundaryMode.FIXED)
                plot.domainStepMode = StepMode.INCREMENT_BY_VAL
                plot.domainStepValue = 1.0
                plot.rangeStepMode = StepMode.INCREMENT_BY_VAL
                plot.rangeStepValue = 0.05
                plot.graph.getLineLabelStyle(XYGraphWidget.Edge.BOTTOM).format = LABEL_FORMATTER
                plot.graph.getLineLabelStyle(XYGraphWidget.Edge.LEFT).format = NUMBER_FORMATTER
                plot.legend.setTableModel(DynamicTableModel(4, 2, TableOrder.ROW_MAJOR))
            }
        }
        plot.redraw()
    }

    fun setParameters(days: Int, idxToDate: Map<Int, Date>) {
        mDays = days;
        mIdxToDate = idxToDate
    }
}
