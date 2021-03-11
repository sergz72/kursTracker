package com.sz.kurstracker

import java.util.Date

interface IGraphData {
    val seriesColumnData: String
    val idx: Int
    val date: Date
    val data: Map<String, Any>
    fun getData(dataName: String): Double?
}
