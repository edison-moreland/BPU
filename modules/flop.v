module flop(input  clk,
            input  [3:0] d,
            output [3:0] q);

    always @(posedge clk)
        q <= d;

endmodule