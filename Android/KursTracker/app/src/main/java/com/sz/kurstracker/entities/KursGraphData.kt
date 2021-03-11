package com.sz.kurstracker.entities

import com.sz.kurstracker.IGraphData
import java.util.*

class KursGraphData(override val seriesColumnData: String, override val idx: Int,
                    override val date: Date) : IGraphData {
    override val data = HashMap<String, Any>()

    override fun getData(dataName: String): Double? {
        return data[dataName] as Double?
    }

    fun setData(dataName: String, value: Double) {
        data[dataName] = value
    }
}