module inc_dec_rst #(
    parameter N = 3
) (
    input  logic         inc, dec, rst,
    input  logic [N-1:0] din,
    output logic [N-1:0] dout,
);
    logic [N-1:0] dinInc;
    logic overflowInc;
    logic [N-1:0] dinDec;
    logic overflowDec;


    always_comb begin
        dout = din;
        {overflowInc, dinInc} = din + 1;
        {overflowDec, dinDec} = din - 1;

        if (rst) dout = 0;
        else if (inc) begin
            if (overflowInc)
                dout = (2**N)-1;
            else
                dout = dinInc;
        end
        else if (dec) begin
            if (overflowDec)
                dout = 0;
            else
                dout = dinDec;
        end
    end
endmodule