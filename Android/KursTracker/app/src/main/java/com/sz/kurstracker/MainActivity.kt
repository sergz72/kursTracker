package com.sz.kurstracker

import android.app.Activity
import android.app.AlertDialog
import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.view.Menu
import android.view.MenuItem
import androidx.appcompat.app.AppCompatActivity
import com.androidplot.xy.XYPlot
import com.google.gson.Gson
import com.google.gson.GsonBuilder
import com.google.gson.reflect.TypeToken
import com.sz.kurstracker.entities.KursData
import com.sz.kurstracker.entities.KursGraphData
import java.util.*
import kotlin.collections.HashMap

class MainActivity : AppCompatActivity() {
    companion object {
        const val PREFS_NAME = "KursTracker"
        private const val SETTINGS = 1

        fun alert(activity: Activity, message: String?) {
            val builder = AlertDialog.Builder(activity)
            builder.setMessage(message)
                    .setTitle("Error")
            val dialog = builder.create()
            dialog.show()
        }
    }

    private val mGson: Gson
    private val mHandler = Handler(Looper.getMainLooper())
    private var mDays = 2
    private val mIdxToDate: MutableMap<Int, Date> = mutableMapOf()

    init {
        val gsonBuilder = GsonBuilder()
        mGson = gsonBuilder.create()
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        setSupportActionBar(findViewById(R.id.toolbar))

        KursService.setupKeyAndContext(resources.openRawResource(R.raw.key), this)

        try {
            updateServer()
        } catch (e: Exception) {
            alert(this, e.message)
        }
    }

    private fun updateServer() {
        val settings = getSharedPreferences(PREFS_NAME, 0)
        val name = settings.getString("server_name", "localhost")!!

        KursService.setupServer(name, 60000)

        refresh()
    }

    private fun refresh() {
        KursService<Array<KursData>>()
                .doInBackground(this, "GET /kurs_data?period=${mDays}", object : KursService.Callback<Array<KursData>> {
                    override fun deserialize(response: String): Array<KursData> {
                        return mGson.fromJson(response, object : TypeToken<Array<KursData>>() {}.type)
                    }

                    override fun isString(): Boolean {
                        return false
                    }

                    override fun onResponse(response: Array<KursData>) {
                        mHandler.post { showResults(response) }
                    }

                    override fun onFailure(activity: Activity, t: Throwable?, response: String?) {
                        mHandler.post {
                            if (t != null) {
                                alert(activity, t.message)
                            } else {
                                alert(activity, response)
                            }
                        }
                    }
                })
    }

    private fun showResults(rawResponse: Array<KursData>) {
        val response = convertResponse(rawResponse)
        var plot = findViewById<XYPlot>(R.id.plot_usd)
        Graph.setParameters(mDays, mIdxToDate)
        Graph.buildGraph(plot, response, "USD")
        plot = findViewById(R.id.plot_eur)
        Graph.buildGraph(plot, response, "EUR")
    }

    private fun convertResponse(rawResponse: Array<KursData>): List<IGraphData> {
        val result : MutableList<IGraphData> = mutableListOf()
        convertResponse(rawResponse, result, "MBB", "USD", "MB", true)
        convertResponse(rawResponse, result, "MBS", "USD", "MB", false)
        convertResponse(rawResponse, result, "MBB", "EUR", "MB", true)
        convertResponse(rawResponse, result, "MBS", "EUR", "MB", false)

        convertResponse(rawResponse, result, "PBB", "USD", "Privat", true)
        convertResponse(rawResponse, result, "PBS", "USD", "Privat", false)
        convertResponse(rawResponse, result, "PBB", "EUR", "Privat", true)
        convertResponse(rawResponse, result, "PBS", "EUR", "Privat", false)

        convertResponse(rawResponse, result, "ALB", "USD", "Alfa", true)
        convertResponse(rawResponse, result, "ALS", "USD", "Alfa", false)
        convertResponse(rawResponse, result, "ALB", "EUR", "Alfa", true)
        convertResponse(rawResponse, result, "ALS", "EUR", "Alfa", false)

        convertResponse(rawResponse, result, "MOB", "USD", "Mono", true)
        convertResponse(rawResponse, result, "MOS", "USD", "Mono", false)
        convertResponse(rawResponse, result, "MOB", "EUR", "Mono", true)
        convertResponse(rawResponse, result, "MOS", "EUR", "Mono", false)

        return result
    }

    private fun convertResponse(rawResponse: Array<KursData>,
                                result: MutableList<IGraphData>,
                                dataName: String,
                                currency: String,
                                bankName: String,
                                isBuy: Boolean) {
        val map = rawResponse.groupBy { it.date }.toMap()
        var idx = 0
        map.keys.sorted().forEach{ k ->
            val d = Date(k * 1000)
            mIdxToDate[idx] = d
            val kursGraphData = KursGraphData(dataName, idx++, d)
            map[k]!!.forEach{ vv ->
                if (vv.bankName == bankName) {
                    if (currency == "USD") {
                        if (isBuy) {
                            kursGraphData.setData(currency, vv.rateBuyUSD)
                        } else {
                            kursGraphData.setData(currency, vv.rateSellUSD)
                        }
                    } else {
                        if (isBuy) {
                            kursGraphData.setData(currency, vv.rateBuyEUR)
                        } else {
                            kursGraphData.setData(currency, vv.rateSellEUR)
                        }
                    }
                }
            }
            result.add(kursGraphData)
        }
    }

    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        menuInflater.inflate(R.menu.menu_main, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            R.id.action_settings -> {
                val intent = Intent(this, SettingsActivity::class.java)
                startActivityForResult(intent, SETTINGS)
            }
            R.id.action_oneday -> {
                mDays = 1
                refresh()
            }
            R.id.action_twodays -> {
                mDays = 2
                refresh()
            }
            R.id.action_oneweek -> {
                mDays = 7
                refresh()
            }
            R.id.action_twoweeks -> {
                mDays = 14
                refresh()
            }
            R.id.action_onemonth -> {
                mDays = 30
                refresh()
            }
            else -> return super.onOptionsItemSelected(item)
        }
        return true
    }

    override fun onActivityResult(requestCodeIn: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCodeIn, resultCode, data)
        var requestCode = requestCodeIn
        requestCode = requestCode and 0xFFFF
        if (requestCode == SETTINGS && resultCode == Activity.RESULT_OK) {
            try {
                updateServer()
            } catch (e: Exception) {
                alert(this, e.message)
            }
        }
    }
}

